package handlers

import "net/http"

func Set_Handlers() {

  http.HandleFunc( "/", handler_main )

}
