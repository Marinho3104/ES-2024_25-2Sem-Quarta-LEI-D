package app

import (
	"fmt"
	"github.com/dominikbraun/graph"
	"time"
)

// globalOwnerGraph is a global graph where vertices are properties and edges represent relationships.
var globalOwnerGraph graph.Graph[int, Property]

// GetOwnerGraph returns the global owner-based graph, initializing it if it has not been created yet.
func GetOwnerGraph() graph.Graph[int, Property] {
	if globalOwnerGraph == nil {
		fmt.Println("Graph has not been created yet, creating now.")
		createOwnerGraph()
	}
	return globalOwnerGraph
}

// SetOwnerGraph sets the globalOwnerGraph variable to the provided graph.
func SetOwnerGraph(a graph.Graph[int, Property]) {
	globalOwnerGraph = a
}

// createOwnerGraph generates a new owner-based graph from the global property graph, linking properties by shared ownership.
// It initializes a graph where vertices represent properties and edges represent ownership relationships between them.
// The function iterates through the property graph's vertices and builds the corresponding owner graph.
func createOwnerGraph() {
	ownerHash := func(p Property) int {
		return p.Owner
	}
	globalOwnerGraph = graph.New(ownerHash)
	propertyGraph := GetGraph()
	adjMap, _ := propertyGraph.AdjacencyMap()
	start := time.Now()
	for vertex := range adjMap {
		property, _ := propertyGraph.Vertex(vertex)
		err := globalOwnerGraph.AddVertex(property)
		if err != nil {
			continue
		}
	}
	edges, _ := propertyGraph.Edges()
	for _, ed := range edges {
		property1, _ := propertyGraph.Vertex(ed.Source)
		property2, _ := propertyGraph.Vertex(ed.Target)
		newIDs := map[int]struct{}{
			property1.Id: {},
			property2.Id: {},
		}

		err2 := globalOwnerGraph.AddEdge(property1.Owner, property2.Owner, graph.EdgeData(newIDs))
		if err2 != nil {
			existingEdge, _ := globalOwnerGraph.Edge(property1.Owner, property2.Owner)
			a := existingEdge.Properties.Data
			m := a.(map[int]struct{})
			for id := range newIDs {
				m[id] = struct{}{}
			}
			_ = globalOwnerGraph.UpdateEdge(property1.Owner, property2.Owner, graph.EdgeData(m))
			continue
		}

	}
	end := time.Now()
	size, _ := globalOwnerGraph.Size()
	order, _ := globalOwnerGraph.Order()
	fmt.Println("Finished OwnerGraph with ", size, " edges ", order, " vertices")
	fmt.Println("In: ", end.Sub(start).Seconds(), "Seconds")
}
