package app

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/dominikbraun/graph"
	"github.com/dominikbraun/graph/draw"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/geojson"
	"github.com/twpayne/go-geom/encoding/wkt"
)

var propertyHash map[string][]Property = make(map[string][]Property)

var propertyHashCode = func(p Property) int {
	return p.id
}

var g = graph.New(propertyHashCode)

func readFile(filename string) ([][]string, error) {

	fmt.Println("Loading file...")

	file, err := os.Open(filename)

	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';'

	var records [][]string

	for {
		record, err := reader.Read()

		if err == io.EOF {
			break
		}

		records = append(records, record)
	}

	return records, nil
}

func createPropertyList() []Property {

	var propertyList []Property = make([]Property, 0)

	data, err := readFile("../../assets/madeira.csv")
	if err != nil {
		fmt.Println("Erro trying reading the file")
		fmt.Println(err)
		return nil
	}

	for i, line := range data {
		// Ignoring first line (header)
		if i > 0 {
			var record Property

			for j, field := range line {
				switch j {
				case 0:
					record.id, _ = strconv.Atoi(field)

				case 4:
					var value, _ = strconv.ParseFloat(field, 32)
					record.shapeArea = float32(value)

				case 5:
					convertedField, err := wkt.Unmarshal(field)

					if err != nil || strings.Contains(field, "EMPTY") {
						break
					}

					record.geometry =
						*geom.NewMultiPolygonFlat(convertedField.Layout(), convertedField.FlatCoords(), convertedField.Endss())

					geojson.Marshal(record.geometry.Clone())
				case 6:
					record.owner, _ = strconv.Atoi(field)
				case 7:
					record.freguesia = field
				case 8:
					record.municipio = field
				default:
					panic("Unreconized field")
				}
			}

			getCoods(line[5], record)

			// Do not add wrong data
			if record.shapeArea == 0 || record.geometry.Bounds().IsEmpty() {
				continue
			}

			propertyList = append(propertyList, record)
		}
	}

	// for _, property := range propertyHash["299218.5203999998 3623637.4791"] {
	// 	fmt.Println(property.id)
	// }

	return propertyList
}

func getCoods(mp string, property Property) {
	result := strings.Replace(mp, "MULTIPOLYGON (((", "", 1)
	result = strings.Replace(result, ")))", "", 1)

	results := strings.SplitSeq(result, ", ")

	for result := range results {
		if propertyHash[result] == nil {
			propertyHash[result] = []Property{property}

		} else if !contains(propertyHash[result], property) {
			propertyHash[result] = append(propertyHash[result], property)
		}
	}
}

func CreateGraph() {

	propertyList := createPropertyList()

	fmt.Println("Creating the graph...")

	for _, property := range propertyList {
		g.AddVertex(property)
	}

	createAdges()

	// size, _ := g.Size()
	// fmt.Println("Fineshed with ", size, " edges")
	// fmt.Println("In ", end.Sub(start).Seconds())

	file, _ := os.Create("../../assets/graph.gv")
	draw.DOT(g, file)
}

func createAdges() {
	for _, propertyHashList := range propertyHash {
		for i := range propertyHashList {
			for j := i + 1; j < len(propertyHashList); j++ {
				g.AddEdge(propertyHashList[i].id, propertyHashList[j].id)
			}
		}
	}
}

func contains(s []Property, e Property) bool {
	for _, a := range s {
		if a.id == e.id {
			return true
		}
	}
	return false
}
