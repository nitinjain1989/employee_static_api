package routes

import (
	"net/http"
	"static-api/config"

	"static-api/handlers"
	"static-api/repositories"
	"static-api/services"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes() http.Handler {
	pgPool := config.InitDB()
	// Dependency Injection
	repo := repositories.NewEmployeeRepository(pgPool)
	service := services.NewEmployeeService(repo)
	handler := handlers.NewEmployeeHandler(service)

	r := gin.Default()

	api := r.Group("/api")

	// ✅ Employees
	api.GET("/employees", handler.GetEmployees)
	api.GET("/employees/:id", handler.GetEmplyoeeByID)

	api.POST("/employees", handler.CreateEmployee)

	// ✅ With ID

	api.PATCH("/employees/:id", handler.UpdateEmployee)

	// ✅ Filters
	api.GET("/employees/filters", handler.GetEmployeeFilters)

	api.DELETE("/employees/:id", handler.DeleteEmployee)

	return r
}
