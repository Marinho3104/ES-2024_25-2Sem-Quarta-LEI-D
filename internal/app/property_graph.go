package app

import (
	"fmt"
	"time"

	"github.com/dhconnelly/rtreego"
	"github.com/dominikbraun/graph"
)

// globalGraph is a global variable storing a graph of Property vertices and edges represent property neighbor relationship.
var globalGraph graph.Graph[int, Property]

// GetGraph returns the global property-based graph, initializing it if it has not been created.
func GetGraph() graph.Graph[int, Property] {
	if globalGraph == nil {
		fmt.Println("Graph has not been created yet, creating now.")
		CreateGraph()
	}
	return globalGraph
}

// SetGraph sets the global graph instance with the provided graph of Property vertices and edges.
func SetGraph(a graph.Graph[int, Property]) {
	globalGraph = a
}

// CreateRTree constructs and returns an R-tree populated with the given list of Properties.
// It initializes the R-tree with default parameters and inserts each Property instance into the tree.
func CreateRTree(propertyList []Property) *rtreego.Rtree {

	rTree := rtreego.NewTree(2, 8, 16)

	for _, property := range propertyList {
		rTree.Insert(&property)
	}

	return rTree

}

// CreateGraph initializes and builds a global graph of Properties and their neighbor relationships using spatial data.
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
	propertyList := getPropertiesList()

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

	fmt.Println("Finished with ", size, " edges")
	fmt.Println("In: ", end.Sub(start).Seconds(), "Seconds")
	uniqueOwners := CountUniqueOwners(propertyList)
	fmt.Println("Number of unique owners:", uniqueOwners)

}

// checkIfNeigbours identifies and adds edges between neighboring Properties in the global graph based on spatial relationships.
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

func CountUniqueOwners(propertyList []Property) int {
	ownerSet := make(map[int]struct{})

	for _, property := range propertyList {
		ownerSet[property.Owner] = struct{}{}
	}

	return len(ownerSet)
}
