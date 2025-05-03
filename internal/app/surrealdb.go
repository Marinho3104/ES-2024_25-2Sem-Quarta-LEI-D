package app

import (
	"fmt"

	"github.com/surrealdb/surrealdb.go"
)

var dbClient *surrealdb.DB

func init() {
	var err error
	// 1. Connect over WebSocket
	dbClient, err = surrealdb.New("ws://127.0.0.1:8000")
	if err != nil {
		panic(fmt.Errorf("failed to connect to SurrealDB: %w", err))
	}

	// 2. Sign in
	auth := &surrealdb.Auth{
		Username: "root",
		Password: "root",
	}
	if _, err = dbClient.SignIn(auth); err != nil {
		panic(fmt.Errorf("failed to authenticate to SurrealDB: %w", err))
	}
	// 3. Select namespace and database
	if err = dbClient.Use("testNS", "testDB"); err != nil {
		panic(fmt.Errorf("failed to switch NS/DB: %w", err))
	}
}
