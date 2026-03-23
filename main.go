package main

import (
	"log"
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

// @title Employee API
// @version 1.0
// @description Employee management service
// @host localhost:8080
// @BasePath /api
func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// API routes
	router := routes.RegisterRoutes()

	/*	// ✅ Serve public folder
		fs := http.FileServer(http.Dir("./public"))
		http.Handle("/", fs)
		http.Handle("/api/", router)*/

	log.Println("Server running on :" + port)
	log.Fatal(router.Run(":" + port))

}
