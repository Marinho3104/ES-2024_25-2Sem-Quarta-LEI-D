package app

import (
	"ES-2024_25-2Sem-Quarta-LEI-D/internal/handlers"
	"net/http"
)

func Start() {

  // Sets all handlers for the web app
  handlers.Set_Handlers()

  // Starts the web app
  http.ListenAndServe( ":8080", nil );

}
