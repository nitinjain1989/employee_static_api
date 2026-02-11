package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"static-api/handlers"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// API routes
	http.HandleFunc("/employees", handlers.GetEmployees)
	http.HandleFunc("/employees/create", handlers.CreateEmployee)
	http.HandleFunc("/employees/update", handlers.UpdateEmployee)

	// Absolute path fix (IMPORTANT)
	publicPath, err := filepath.Abs("public")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Serving static files from:", publicPath)

	http.Handle("/", http.FileServer(http.Dir(publicPath)))

	log.Println("Server started on port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
