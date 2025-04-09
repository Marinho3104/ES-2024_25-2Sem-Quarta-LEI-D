package handlers

import "net/http"

func Start() {

	//http.HandleFunc("/", handler_main)
	http.HandleFunc("/api/upload", upload_handler)

	http.HandleFunc("/api/graph", graphdata_handler)

	http.ListenAndServe(":8080", nil)
}
