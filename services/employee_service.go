package services

import (
	"context"
	"fmt"
	"static-api/models"
	"static-api/repositories"
)

type EmployeeService struct {
	repo *repositories.EmployeeRepository
}

func NewEmployeeService(r *repositories.EmployeeRepository) *EmployeeService {
	return &EmployeeService{repo: r}
}

func (s *EmployeeService) GetEmployees(
	ctx context.Context,
	filter models.EmployeeFilter,
) ([]models.Employee, models.Meta, error) {

	return s.repo.FetchEmployees(ctx, filter)
}

/*func (s *EmployeeService) GetEmployees(c *gin.Context) ([]models.Employee, models.Meta, error) {

	path := utils.BuildPath(c)

	resp, err := s.repo.FetchEmployees(path)

	if err != nil {
		return nil, models.Meta{}, err
	}

	meta := utils.ParseContentRange(resp.ContentRange, c.Query("limit"))

	return resp.Data, meta, nil

}*/

func (s *EmployeeService) GetEmployeeByID(id string) (*models.Employee, error) {
	employees, err := s.repo.GetEmployeeByID(context.Background(), id)
	if err != nil {
		return nil, err
	}

	return employees, nil
}

func (s *EmployeeService) CreateEmployee(emp models.Employee) (string, error) {

	return s.repo.CreateEmployeeWithMobiles(context.Background(), emp)
}

func (s *EmployeeService) GetEmployeeFilters() (*models.EmployeeFilters, error) {

	designations, err := s.repo.FetchDistinctField(context.Background(), "designation")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch designations")
	}

	departments, err := s.repo.FetchDistinctField(context.Background(), "department")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch departments")
	}

	statuses := []models.FilterOption{
		{Label: "Active", Value: "active"},
		{Label: "Inactive", Value: "inactive"},
	}

	mobileTypes := []models.FilterOption{
		{Label: "Home", Value: "home"},
		{Label: "Office", Value: "office"},
		{Label: "Other", Value: "other"},
	}

	return &models.EmployeeFilters{
		Designations: designations,
		Departments:  departments,
		Statuses:     statuses,
		MobileTypes:  mobileTypes,
	}, nil
}

func (s *EmployeeService) UpdateEmployee(id string, emp models.Employee) error {
	return s.repo.UpdateEmployeeWithMobiles(context.Background(), id, emp)
}

func (s *EmployeeService) DeleteEmployee(id string) error {
	return s.repo.DeleteEmployee(context.Background(), id)
}
