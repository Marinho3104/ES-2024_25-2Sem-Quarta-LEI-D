package handlers

import "net/http"

func Set_Handlers() {

  fs := http.FileServer(http.Dir("web"))
  http.Handle("/", http.StripPrefix("/", fs))

}
