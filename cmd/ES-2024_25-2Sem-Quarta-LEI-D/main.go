package main

import (
	//"ES-2024_25-2Sem-Quarta-LEI-D/internal/app"
	"ES-2024_25-2Sem-Quarta-LEI-D/internal/app"
	"ES-2024_25-2Sem-Quarta-LEI-D/internal/handlers"
)

func main() {
	app.GetGraph()
	handlers.Start()
}
