package app

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/dominikbraun/graph"
	"github.com/dominikbraun/graph/draw"

	// "github.com/dominikbraun/graph/draw"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/geojson"
	"github.com/twpayne/go-geom/encoding/wkt"
)

var globalGraph graph.Graph[int, Property]

func readFile() ([][]string, error) {

	fmt.Println("Loading file...")

	file, err := os.Open("../../assets/madeira.csv")

	if err != nil {
		return nil, err
	}
	defer file.Close()

	fmt.Println("File loaded successfuly!!")

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

	var propertyList []Property

	data, err := readFile()
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
				case 1:
				case 2:
				case 3:
				case 4:
					var value, _ = strconv.ParseFloat(field, 32)
					record.shapeArea = float32(value)
				case 5:
					convertedField, err := wkt.Unmarshal(field)

					if err != nil {
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
				case 9:
				default:
					panic("Unreconized field")
				}
			}

			// Do not add wrong data
			if record.shapeArea == 0 || record.geometry.Bounds().IsEmpty() {
				continue
			}

			propertyList = append(propertyList, record)
		}
	}

	return propertyList
}

func CreateGraph() {

	if globalGraph != nil {
		fmt.Println("Graph already created, returning the existing graph.")
		return
	}

	propertyHash := func(p Property) int {
		return p.id
	}

	globalGraph = graph.New(propertyHash)
	propertyList := createPropertyList()

	fmt.Println("Creating the graph")

	for _, property := range propertyList {
		globalGraph.AddVertex(property)
	}

	fmt.Println("Creating edges")

	propertyListLen := len(propertyList)

	start := time.Now()
	for i := 0; i < propertyListLen; i++ {
		currentProperty := propertyList[i]

		for j := i + 1; j < propertyListLen; j++ {
			cmpProperty := propertyList[j]
			// Check if bounding boxes intersect
			if bboxOverlap(currentProperty.geometry.Bounds(), cmpProperty.geometry.Bounds()) {

				// fmt.Printf("%d --- %d\n", propertyList[i].id, propertyList[j].id)
				err := globalGraph.AddEdge(currentProperty.id, cmpProperty.id)
				if err != nil {
					fmt.Println(err)
				}

			}

		}
	}
	end := time.Now()

	size, _ := globalGraph.Size()

	fmt.Println("Fineshed with ", size, " edges")
	fmt.Println("In: ", end.Sub(start).Minutes())

	file, _ := os.Create("../../assets/graph.gv")
	draw.DOT(globalGraph, file)

}

func bboxOverlap(a, b *geom.Bounds) bool {
	return !(a.Max(0) < b.Min(0) || a.Min(0) > b.Max(0) ||
		a.Max(1) < b.Min(1) || a.Min(1) > b.Max(1))
}

func GetGraph() graph.Graph[int, Property] {
	if globalGraph == nil {
		fmt.Println("Graph has not been created yet, creating now.")
		CreateGraph()
	}
	return globalGraph
}

func SetGraph(a graph.Graph[int, Property]) {
	globalGraph = a
}
