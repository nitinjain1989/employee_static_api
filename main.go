package main

import (
	"log"
	"net/http"
	"os"
	"static-api/routes"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}
}
func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	/*// API routes
	router := routes.RegisterRoutes()
	//http.HandleFunc("/employees", handlers.GetEmployees)
	//http.HandleFunc("/employees/get", handlers.GetEmployeeByID)
	http.HandleFunc("/employees/create", handlers.CreateEmployee)
	http.HandleFunc("/employees/update", handlers.UpdateEmployee)
	http.HandleFunc("/employees/filters", handlers.GetEmployeeFilters)

	// Absolute path fix (IMPORTANT)
	publicPath, err := filepath.Abs("public")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Serving static files from:", publicPath)

	http.Handle("/", http.FileServer(http.Dir(publicPath)))

	log.Println("Server started on port", port)
	log.Fatal(http.ListenAndServe(":"+port, router))*/

	//mux := http.NewServeMux()

	// API routes
	router := routes.RegisterRoutes()

	// ✅ Serve public folder
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)
	http.Handle("/api/", router)

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
