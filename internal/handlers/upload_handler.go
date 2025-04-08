package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func upload_handler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/api/upload" {
		http.NotFound(w, r)
		return
	}

	fmt.Println("Upload endpoint")

	file, header, err := r.FormFile("file")

	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}
	defer file.Close()

	dest, err := os.Create("../../assets/" + header.Filename)
	if err != nil {
		fmt.Println(err)
	}
	defer dest.Close()

	_, err = io.Copy(dest, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "File %s uploaded successfully", header.Filename)
}
