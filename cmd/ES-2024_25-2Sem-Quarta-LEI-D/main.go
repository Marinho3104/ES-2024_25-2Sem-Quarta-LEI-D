package main

import (
	//"ES-2024_25-2Sem-Quarta-LEI-D/internal/app"
	"ES-2024_25-2Sem-Quarta-LEI-D/internal/app"
	// "ES-2024_25-2Sem-Quarta-LEI-D/internal/handlers"
	//"ES-2024_25-2Sem-Quarta-LEI-D/internal/handlers"
	//"ES-2024_25-2Sem-Quarta-LEI-D/internal/versiondb"
)

// main initializes the program by setting up the required graphs and starting HTTP handlers.
func main() {

	app.Init_prices_by_concelho_and_fregusesia();

	//versiondb.Main()
	// app.GetGraph()
	// app.GetOwnerGraph()
	// handlers.Start()

}
