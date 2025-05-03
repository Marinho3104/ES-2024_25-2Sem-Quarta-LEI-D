package app

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/dhconnelly/rtreego"
	"github.com/dominikbraun/graph"
	"github.com/dominikbraun/graph/draw"
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
		if i > 0 {
			prop, err := createProperty(line)
			if err != nil {
				continue
			}
			propertyList = append(propertyList, *prop)
		}

	}
	return propertyList
}

func CreateRTree(propertyList []Property) *rtreego.Rtree {

	rTree := rtreego.NewTree(2, 25, 50)

	for _, property := range propertyList {
		rTree.Insert(&property)
	}

	return rTree

}

func CreateGraph() {

	if globalGraph != nil {
		fmt.Println("Graph already created, returning the existing graph.")
		return
	}

	propertyHash := func(p Property) int {
		return p.Id
	}

	fmt.Println("Creating the graph")
	globalGraph = graph.New(propertyHash)
	propertyList := createPropertyList()

	fmt.Println("Creating the RTree")
	rTree := CreateRTree(propertyList)

	fmt.Println("Creating the vertex")
	for _, property := range propertyList {
		globalGraph.AddVertex(property)
	}

	fmt.Println("Creating edges")

	start := time.Now()

	for _, currentProperty := range propertyList {

		results := rTree.SearchIntersect(*currentProperty.Rect)
		checkIfNeigbours(results)

	}

	end := time.Now()

	size, _ := globalGraph.Size()

	fmt.Println("Fineshed with ", size, " edges")
	fmt.Println("In: ", end.Sub(start).Seconds())

	file, _ := os.Create("../../assets/graph.gv")
	draw.DOT(globalGraph, file)

}

func checkIfNeigbours(potentialNeighbors []rtreego.Spatial) {
	for i := 0; i < len(potentialNeighbors); i++ {
		property1 := potentialNeighbors[i].(*Property)
		for j := i + 1; j < len(potentialNeighbors); j++ {

			property2 := potentialNeighbors[j].(*Property)
			if areMultiPolygonsNeighbors(&property1.Geometry, &property2.Geometry) {
				globalGraph.AddEdge(property1.Id, property2.Id)
			}

		}
	}
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

