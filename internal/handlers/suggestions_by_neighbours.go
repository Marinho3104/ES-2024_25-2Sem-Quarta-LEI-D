package handlers

import (
	"ES-2024_25-2Sem-Quarta-LEI-D/internal/app"
	"encoding/json"
	"net/http"
)

func suggestions_by_neighbours_handler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode( app.Suggestions_point_6 )


}
