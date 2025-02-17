package main

import (
	"Go-Check24/database"
	"Go-Check24/handlers"
	"Go-Check24/router"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	measurementDB, err := database.InitDB()
	if err != nil {
		//no sense continuing if can't open DB
		log.Fatal("Database connection/creation failed:", err)
	}
	defer measurementDB.Close()

	measurementHandler := handlers.NewHandler(measurementDB)
	r := gin.Default()

	router.SetupRoutes(r, measurementHandler)

	err = r.Run(":8080")
	if err != nil {
		measurementDB.Close()             // Ensure proper cleanup
		log.Fatal("Server failed: ", err) // Exit the program
	}
}
