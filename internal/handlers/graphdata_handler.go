package handlers

import (
	"ES-2024_25-2Sem-Quarta-LEI-D/internal/app"
	"encoding/json"
	"net/http"
	"strconv"
)

type SigmaNode struct {
	ID    string  `json:"id"`
	Label string  `json:"label,omitempty"`
	X     float64 `json:"x,omitempty"`
	Y     float64 `json:"y,omitempty"`
}

type SigmaEdge struct {
	ID     string `json:"id"`
	Source string `json:"source"`
	Target string `json:"target"`
	Label  string `json:"label,omitempty"`
}

type SigmaGraphResponse struct {
	Nodes []SigmaNode `json:"nodes"`
	Edges []SigmaEdge `json:"edges"`
}

func graphdata_handler(w http.ResponseWriter, r *http.Request) {
	g := app.GetGraph()

	adjMap, err := g.AdjacencyMap()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	nodes := make([]SigmaNode, 0, len(adjMap))
	for vertex := range adjMap {
		property, err := g.Vertex(vertex)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var longitude, latitude float64

		clonedGeometry := property.Geometry()

		polygon := clonedGeometry.Polygon(0)
		if polygon != nil && polygon.Length() > 0 {
			coords := polygon.Coords()
			if len(coords) > 0 {
				longitude = coords[0][0].X()
				latitude = coords[0][0].Y()
			}
		}
		nodes = append(nodes, SigmaNode{
			ID: strconv.Itoa(vertex),
			X:  longitude,
			Y:  latitude,
		})
	}

	edges := make([]SigmaEdge, 0)
	for source, targets := range adjMap {
		for target := range targets {
			edges = append(edges, SigmaEdge{
				ID:     strconv.Itoa(source) + "-" + strconv.Itoa(target),
				Source: strconv.Itoa(source),
				Target: strconv.Itoa(target),
			})
		}
	}

	response := SigmaGraphResponse{
		Nodes: nodes,
		Edges: edges,
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(response)
}
