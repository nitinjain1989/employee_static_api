package main

import (
	"log"
	"os"
	"static-api/routes"

	"github.com/joho/godotenv"

	docs "static-api/docs"
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

	docs.SwaggerInfo.Title = "Employee API"
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Description = "Employee Management Service"
	docs.SwaggerInfo.Version = "1.0"
	env := os.Getenv("ENV")
	if env == "production" {
		docs.SwaggerInfo.Host = "https://employee-static-api.onrender.com"
		docs.SwaggerInfo.Schemes = []string{"https"}
	} else {
		docs.SwaggerInfo.Host = "localhost:8080"
		docs.SwaggerInfo.Schemes = []string{"http"}
	}

	// API routes
	router := routes.RegisterRoutes()

	log.Println("Server running on :" + port)
	log.Fatal(router.Run(":" + port))

}
