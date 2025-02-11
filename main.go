package main

import (
	db "Go-Check24/database"
	"Go-Check24/server"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	database := db.Database{}
	database.Open()
	defer database.Close()
	//+++++++++++++++++++++++++++++++++++++++++++
	//+++++++++++++++++++++++++++++++++++++++++++
	//BAREBONE SERVER

	http.HandleFunc("/measurement", func(w http.ResponseWriter, r *http.Request) {
		server.Measurement(w, r, database)
	}) // closure on database
	//fmt.Println("Server is running on http://localhost:8080")
	//log.Fatal(http.ListenAndServe(":8080", nil)) //here we set up the server
}
