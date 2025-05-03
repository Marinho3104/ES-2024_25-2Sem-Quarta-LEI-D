package versiondb

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"

	surrealdb "github.com/surrealdb/surrealdb.go"
	"github.com/surrealdb/surrealdb.go/pkg/models"
	"github.com/twpayne/go-geom/encoding/geojson"
)

// dbProperty is what we serialize into SurrealDB (GeoJSON map).
type dbProperty struct {
	Id        *models.RecordID       `json:"id,omitempty"`
	Owner     int                    `json:"owner"`
	ShapeArea float32                `json:"shapeArea"`
	Freguesia string                 `json:"freguesia"`
	Municipio string                 `json:"municipio"`
	Geometry  map[string]interface{} `json:"geometry"`
}

var dbClient *surrealdb.DB

// initSurreal connects, authenticates, and selects NS/DB.
func initSurreal() {
	var err error
	dbClient, err = surrealdb.New("ws://127.0.0.1:8000")
	if err != nil {
		panic(fmt.Errorf("connect error: %w", err))
	}
	if _, err = dbClient.SignIn(&surrealdb.Auth{
		Username: "root",
		Password: "root",
	}); err != nil {
		panic(fmt.Errorf("auth error: %w", err))
	}
	if err = dbClient.Use("testNS", "testDB"); err != nil {
		panic(fmt.Errorf("use error: %w", err))
	}
}

// readCSV loads all records from a semicolon-delimited file.
func readCSV(path string) ([][]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	r := csv.NewReader(f)
	r.Comma = ';'
	var data [][]string
	for {
		rec, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		data = append(data, rec)
	}
	return data, nil
}

// upsertProperties reads the CSV, parses each row, and upserts into DB.
func upsertProperties(csvPath string) error {
	data, err := readCSV(csvPath)
	if err != nil {
		return err
	}
	var props []*dbProperty

	for i, row := range data {
		if i == 0 {
			continue // skip header
		}
		p, err := createProperty(row)
		if err != nil {
			fmt.Printf("skip line %d: %v\n", i, err)
			continue
		}

		geometryJSON, _ := geojson.Marshal(&p.Geometry)

		var geometryMap map[string]interface{}
		json.Unmarshal(geometryJSON, &geometryMap)

		propertyID := models.NewRecordID("property", p.Id)
		dbp := dbProperty{
			Id:        &propertyID,
			Owner:     p.Owner,
			ShapeArea: p.ShapeArea,
			Freguesia: p.Freguesia,
			Municipio: p.Municipio,
			Geometry:  geometryMap,
		}

		props = append(props, &dbp)

	}
	for i, prop := range props {
		_, err = surrealdb.Upsert[dbProperty](dbClient, *prop.Id, prop)
		if err != nil {
			fmt.Printf("upsert failed for %v: %v\n", i, err)
		}
	}
	return nil
}

// setupSchemaAndIndex defines tables and the MTREE index.
func setupSchemaAndIndex(ctx context.Context) error {
	stmts := []string{
		`DEFINE TABLE property SCHEMALESS;`,
	}
	for _, s := range stmts {
		if _, err := surrealdb.Query[any](dbClient, s, map[string]interface{}{}); err != nil {
			return err
		}
	}

	//		`DEFINE TABLE neighbor RELATIONSHIP;`,
	//	`DEFINE INDEX IF NOT EXISTS idx_geom ON property FIELDS geometry USING MTREE TYPE F64;`,

	return nil
}

/*
// buildGraphEdges runs the INTERSECTS+RELATE script and prints the raw result.
func buildGraphEdges(ctx context.Context) error {
	script := `
LET props = (SELECT id, geometry FROM property);
FOR p IN props {
  LET neigh = (
    SELECT id AS to_id
    FROM property
    WHERE id != p.id
      AND geometry INTERSECTS p.geometry
  );
  FOR q IN neigh {
    RELATE property:p.id->neighbor->property:q.to_id;
  }
}`
	res, err := dbClient.Query(ctx, script)

	if err != nil {
		return err
	}
	fmt.Printf("Graph build response: %#v\n", res.Result)
	return nil
}

// queryNeighbors returns all incoming neighbor IDs for a property.
func queryNeighbors(ctx context.Context, propertyID int) ([]int, error) {
	stmt := fmt.Sprintf(`
SELECT neighbor <- neighbor - property:%d AS from
`, propertyID)
	res, err := dbClient.Query(ctx, stmt)
	if err != nil {
		return nil, err
	}
	var arr []struct {
		From int `json:"from"`
	}
	if err := res.All(&arr); err != nil {
		return nil, err
	}
	var ids []int
	for _, e := range arr {
		ids = append(ids, e.From)
	}
	return ids, nil
}



*/

func test() {

	type Person struct {
		ID       *models.RecordID     `json:"id,omitempty"`
		Name     string               `json:"name"`
		Surname  string               `json:"surname"`
		Location models.GeometryPoint `json:"location"`
	}
	personID := models.NewRecordID("persons", "1")

	person2, err := surrealdb.Create[Person](dbClient, personID, Person{
		ID:       &personID,
		Name:     "John",
		Surname:  "Doe",
		Location: models.NewGeometryPoint(-0.11, 22.00),
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created person with a struvt: %+v\n", person2)
}

func Main() {
	ctx := context.Background()
	initSurreal()

	if err := setupSchemaAndIndex(ctx); err != nil {
		panic(err)
	}
	fmt.Println("Schema and MTREE index ready.")

	//test()

	err := upsertProperties("../../assets/madeira.csv")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Upserted  properties.\n")

	/*
		if err := buildGraphEdges(ctx); err != nil {
			panic(err)
		}
		fmt.Println("Graph edges created.")

		// Example query: neighbors of the first property
		sampleID := props[0].Id
		neighbors, err := queryNeighbors(ctx, sampleID)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Neighbors of property %d: %v\n", sampleID, neighbors)
	*/
}
