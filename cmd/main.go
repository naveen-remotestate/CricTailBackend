package main

import (
	"CricTail_Backend/database"
	"CricTail_Backend/server"
	"fmt"
)

func main() {
	dbHost := "localhost"
	dbPort := "5432"
	dbUser := "naveen"
	dbPassword := "naveen"
	dbName := "crictail-db"
	sslMode := database.SSLModeDisable
	serverPort := ":8080"

	err := database.ConnectAndMigrate(
		dbHost,
		dbPort,
		dbName,
		dbUser,
		dbPassword,
		sslMode,
	)
	if err != nil {
		fmt.Printf("Failed to initialise and Migrate database: %v", err)
	}

	router := server.StartServer(serverPort)

	fmt.Println("Server is Running on Port:", serverPort)
	_ = router.Run(serverPort)
}
