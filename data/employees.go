package data

import (
	"fmt"
	"static-api/models"
)

var Employees = generateEmployees()

func generateEmployees() []models.Employee {
	employees := make([]models.Employee, 0, 100)

	for i := 1; i <= 100; i++ {
		emp := models.Employee{
			ID:          fmt.Sprintf("EMP%03d", i),
			Name:        fmt.Sprintf("Employee %d", i),
			Designation: "Software Engineer",
			Department:  "Engineering",
			IsActive:    i%2 == 0,
			ImgURL:      "https://example.com/images/default.png",
			Email:       fmt.Sprintf("employee%d@example.com", i),
			Location: models.Location{
				City:    "Bangalore",
				Country: "India",
			},
			JoiningDate: "2024-01-01",
		}

		employees = append(employees, emp)
	}

	return employees
}
