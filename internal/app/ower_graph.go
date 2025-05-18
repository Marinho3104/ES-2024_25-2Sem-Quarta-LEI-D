package app

import (
	"fmt"
	"github.com/dominikbraun/graph"
	"sort"
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
	globalOwnerGraph = graph.New(ownerHash, graph.Directed())

	propertyGraph := GetGraph()

	// Step 1: Sort vertices before adding
	adjMap, _ := propertyGraph.AdjacencyMap()
	vertexIDs := make([]int, 0, len(adjMap))
	for v := range adjMap {
		vertexIDs = append(vertexIDs, v)
	}
	sort.Ints(vertexIDs)

	for _, vertex := range vertexIDs {
		property, _ := propertyGraph.Vertex(vertex)
		_ = globalOwnerGraph.AddVertex(property)
	}

	// Step 2: Sort edges before processing
	edges, _ := propertyGraph.Edges()
	sort.Slice(edges, func(i, j int) bool {
		return edges[i].Source < edges[j].Source
	})

	// Step 3: Add edges deterministically
	for _, ed := range edges {
		property1, _ := propertyGraph.Vertex(ed.Source)
		property2, _ := propertyGraph.Vertex(ed.Target)

		// Build edge data: only include target property ID
		newIDs := map[int]struct{}{
			property2.Id: {},
		}

		// Add edge: property1.Owner -> property2.Owner
		err := globalOwnerGraph.AddEdge(property1.Owner, property2.Owner, graph.EdgeData(newIDs))
		if err != nil {
			existingEdge, _ := globalOwnerGraph.Edge(property1.Owner, property2.Owner)
			existingData, ok := existingEdge.Properties.Data.(map[int]struct{})
			if !ok {
				existingData = map[int]struct{}{}
			}
			existingData[property2.Id] = struct{}{}
			_ = globalOwnerGraph.UpdateEdge(property1.Owner, property2.Owner, graph.EdgeData(existingData))
		}

		// Repeat for reverse direction: property2.Owner -> property1.Owner
		reverseIDs := map[int]struct{}{
			property1.Id: {},
		}

		err = globalOwnerGraph.AddEdge(property2.Owner, property1.Owner, graph.EdgeData(reverseIDs))
		if err != nil {
			existingEdge, _ := globalOwnerGraph.Edge(property2.Owner, property1.Owner)
			existingData, ok := existingEdge.Properties.Data.(map[int]struct{})
			if !ok {
				existingData = map[int]struct{}{}
			}
			existingData[property1.Id] = struct{}{}
			_ = globalOwnerGraph.UpdateEdge(property2.Owner, property1.Owner, graph.EdgeData(existingData))
		}
	}

	// Step 4: Report stats
	size, _ := globalOwnerGraph.Size()
	order, _ := globalOwnerGraph.Order()
	fmt.Println("Finished OwnerGraph with", size, "edges and", order, "vertices")

}
