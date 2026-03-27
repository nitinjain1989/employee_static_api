package services

import (
	"context"
	"fmt"
	"static-api/dto"
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
	filter dto.EmployeeFilterRequest,
) ([]dto.EmployeeResponse, dto.Meta, error) {

	employees, meta, err := s.repo.FetchEmployees(ctx, filter)
	if err != nil {
		return nil, meta, err
	}

	if len(employees) == 0 {
		return []dto.EmployeeResponse{}, meta, nil
	}

	return employees, meta, nil
}

func (s *EmployeeService) GetEmployeeByID(id string) (*dto.EmployeeResponse, error) {
	employees, err := s.repo.GetEmployeeByID(context.Background(), id)
	if err != nil {
		return nil, err
	}

	return employees, nil
}

func (s *EmployeeService) CreateEmployee(emp dto.CreateEmployeeRequest) (*dto.EmployeeResponse, error) {

	return s.repo.CreateEmployeeWithMobiles(context.Background(), emp)
}

func (s *EmployeeService) GetEmployeeFilters() (*dto.EmployeeFilters, error) {

	designations, err := s.repo.FetchDistinctField(context.Background(), "designation")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch designations")
	}

	departments, err := s.repo.FetchDistinctField(context.Background(), "department")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch departments")
	}

	statuses := []dto.FilterOption{
		{Label: "Active", Value: "active"},
		{Label: "Inactive", Value: "inactive"},
	}

	mobileTypes := []dto.FilterOption{
		{Label: "Home", Value: "home"},
		{Label: "Office", Value: "office"},
		{Label: "Other", Value: "other"},
	}

	return &dto.EmployeeFilters{
		Designations: designations,
		Departments:  departments,
		Statuses:     statuses,
		MobileTypes:  mobileTypes,
	}, nil
}

func (s *EmployeeService) UpdateEmployee(id string, emp dto.UpdateEmployeeRequest) (*dto.EmployeeResponse, error) {
	return s.repo.UpdateEmployeeWithMobiles(context.Background(), id, emp)
}

func (s *EmployeeService) DeleteEmployee(id string) error {
	return s.repo.DeleteEmployee(context.Background(), id)
}
