package main

import (
	"log"
	"net/http"
	"static-api/handlers"
)

func main() {

	http.HandleFunc("/employees", handlers.GetEmployees)
	log.Println("Server started on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
