package handlers

import (
	"ES-2024_25-2Sem-Quarta-LEI-D/internal/app"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

func Area_average(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Parse query parameters
	name := r.URL.Query().Get("name")
	levelStr := r.URL.Query().Get("type")

	if name == "" || levelStr == "" {
		http.Error(w, `{"error": "missing 'name' or 'type' parameter"}`, http.StatusBadRequest)
		return
	}

	name = strings.TrimSpace(name)
	level, err := strconv.Atoi(levelStr)
	if err != nil || level < 0 || level > 2 {
		http.Error(w, `{"error": "invalid 'type' parameter. Must be 0 (Distrito), 1 (Municipio), or 2 (Freguesia)"}`, http.StatusBadRequest)
		return
	}

	// Calculate the average area
	avg, err := app.CalcOfArea(name, level)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusNotFound)
		return
	}

	// Return the result
	response := map[string]interface{}{
		"name":         name,
		"type":         level,
		"average_area": avg,
	}

	json.NewEncoder(w).Encode(response)
}
