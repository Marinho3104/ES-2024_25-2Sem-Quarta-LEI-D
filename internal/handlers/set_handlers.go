package handlers

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
)

func Start() {

	fmt.Println("Starting backend...")

	http.HandleFunc("/", handler_main)
	http.HandleFunc("/api/upload", upload_handler)
	http.HandleFunc("/api/prop", property_handler)
	http.HandleFunc("/api/graph", graphdata_handler)

	startFrontEnd();

	http.ListenAndServe(":8080", nil);
}

func startFrontEnd() {
	fmt.Println("Starting frontend...");

	// Install all npm dependecies
	runNpmCommand("i")

	// Start React website
	runNpmCommand("start")
}

func runNpmCommand(arg ...string) {
	const reactAppPath = "../../es-project-react-app/";

	cmd := exec.Command("npm", arg...)
	cmd.Dir = reactAppPath;

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error running npm command: %v\n", err);
		os.Exit(-1);
	}
}
