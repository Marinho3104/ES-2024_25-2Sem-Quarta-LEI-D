package handlers

import (
	"ES-2024_25-2Sem-Quarta-LEI-D/internal/app"
	"encoding/json"
	"net/http"
	"strconv"

	//"github.com/paulmach/go.geojson"
	"github.com/twpayne/go-geom/encoding/geojson"
)

func property_handler(w http.ResponseWriter, r *http.Request) {
	fc := &geojson.FeatureCollection{
		Features: []*geojson.Feature{},
	}

	g := app.GetGraph()

	adjMap, err := g.AdjacencyMap()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for vertex := range adjMap {
		p, err := g.Vertex(vertex)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		feature := &geojson.Feature{
			Geometry: &p.Geometry,
			Properties: map[string]interface{}{
				"id":        p.Id,
				"owner":     p.Owner,
				"area":      p.ShapeArea,
				"freguesia": p.Freguesia,
				"municipio": p.Municipio,
			},
		}

		feature.Properties["id"] = p.Id
		feature.Properties["owner"] = p.Owner
		feature.Properties["area"] = p.ShapeArea
		feature.Properties["freguesia"] = p.Freguesia
		feature.Properties["municipio"] = p.Municipio

		fc.Features = append(fc.Features, feature)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(fc)
}

func property_neighbour_handler(w http.ResponseWriter, r *http.Request) {
	// Get the "id" parameter from the query string
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Missing 'id' query parameter", http.StatusBadRequest)
		return
	}

	// Convert the id to int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "'id' must be an integer", http.StatusBadRequest)
		return
	}

	// Get neighbors
	neighbors, ok := app.GetIDPropertyNeighbour(id)
	if !ok {
		http.Error(w, "Could not find neighbors for the given property ID", http.StatusNotFound)
		return
	}

	// Set response headers and encode the result to JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(neighbors); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
