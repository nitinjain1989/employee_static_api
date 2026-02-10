package main

import (
	"log"
	"net/http"
	"os"
	"static-api/handlers"
)

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/employees", handlers.GetEmployees)
	log.Println("Server started on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
