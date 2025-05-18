package handlers

import "net/http"

func Start() {

	//http.HandleFunc("/", handler_main)
	http.HandleFunc("/api/upload", upload_handler)

	http.HandleFunc("/api/prop", property_handler)

	http.HandleFunc("/api/graph", graphdata_handler)

	http.HandleFunc("/api/owner", grapownerhdata_handler)

	http.HandleFunc("/api/adm-area", administrative_area_handler)

	http.HandleFunc("/api/prop-neighbour", property_neighbour_handler)

	http.HandleFunc("/api/area-average", Area_average)

	http.HandleFunc("/api/suggestions_by_neighbours", suggestions_by_neighbours_handler)

	http.ListenAndServe(":8080", nil)

}
