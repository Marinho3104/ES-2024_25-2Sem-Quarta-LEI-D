package handlers

import (
	"net/http"
)

func handler_main( w http.ResponseWriter, r *http.Request ) {
    if r.URL.Path != "/" {
        http.NotFound(w, r)
        return
    }

  http.ServeFile(w, r, "./web/main/main.html")

}
