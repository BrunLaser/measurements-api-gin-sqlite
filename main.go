package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Measurement struct {
	ID int `json:"id"`
	//Quantity  string  `json:"quantity"`
	//Timestamp string  `json:"timestamp"`
	Value float64 `json:"value"`
}

var m = []Measurement{
	{ID: 1111, Value: 1.111},
	{ID: 2222, Value: 2.222},
}

func messpunkte(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json") //set the header
		json.NewEncoder(w).Encode(m)                       //write json data
	case http.MethodPost:
		//check for json header
		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "Kein JSON Header", http.StatusBadRequest)
			return
		}

		newM := &Measurement{}

		if err := json.NewDecoder(r.Body).Decode(newM); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		newM.ID = len(m) + 1
		m = append(m, *newM)

		//location header
		location := fmt.Sprintf("/messpunkte/%d", newM.ID)
		w.Header().Set("Location", location)
		//Anwort senden
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(*newM)

	default:
		http.Error(w, "Unerlaubte Methode", http.StatusMethodNotAllowed)
	}

}

func main() {
	//m := &Measurement{
	//ID: 1234, Quantity: "velocity", Timestamp: "Monday", Value: 3.1415}

	http.HandleFunc("/messpunkte", messpunkte) // here we register the handle func
	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil)) //here we set up the server
}
