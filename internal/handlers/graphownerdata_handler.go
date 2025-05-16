package handlers

import (
	"ES-2024_25-2Sem-Quarta-LEI-D/internal/app"
	"encoding/json"
	"net/http"
	"strconv"
)

type SigmaNodeOwner struct {
	ID    string  `json:"id"`
	Label string  `json:"label,omitempty"`
	X     float64 `json:"x,omitempty"`
	Y     float64 `json:"y,omitempty"`
}

type SigmaEdgeOwner struct {
	ID     string `json:"id"`
	Source string `json:"source"`
	Target string `json:"target"`
	Label  string `json:"label,omitempty"`
}

type SigmaGraphOwnerResponse struct {
	Nodes []SigmaNode `json:"nodes"`
	Edges []SigmaEdge `json:"edges"`
}

func grapownerhdata_handler(w http.ResponseWriter, r *http.Request) {
	g := app.GetOwnerGraph()

	adjMap, err := g.AdjacencyMap()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	nodes := make([]SigmaNode, 0, len(adjMap))
	for vertex := range adjMap {
		nodes = append(nodes, SigmaNode{
			ID: strconv.Itoa(vertex),
			X:  0,
			Y:  0,
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
