package routes

import (
	"static-api/config"

	"static-api/handlers"
	"static-api/repositories"
	"static-api/services"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "static-api/docs" // 👈 IMPORTANT
)

func RegisterRoutes() *gin.Engine {
	pgPool := config.InitDB()
	// Dependency Injection
	repo := repositories.NewEmployeeRepository(pgPool)
	service := services.NewEmployeeService(repo)
	handler := handlers.NewEmployeeHandler(service)

	r := gin.Default()

	api := r.Group("/api")

	r.GET("/", func(c *gin.Context) {
		c.File("./public/index.html")
	})

	r.GET("/add", func(c *gin.Context) {
		c.File("./public/add.html")
	})

	r.Static("/js", "./public/js")

	// ✅ Employees
	api.GET("/employees", handler.GetEmployees)
	api.GET("/employees/:id", handler.GetEmplyoeeByID)

	api.POST("/employees", handler.CreateEmployee)

	// ✅ With ID

	api.PATCH("/employees/:id", handler.UpdateEmployee)

	// ✅ Filters
	api.GET("/employees/filters", handler.GetEmployeeFilters)

	api.DELETE("/employees/:id", handler.DeleteEmployee)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
