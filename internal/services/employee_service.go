package services

import (
	"context"

	"app/internal/models"
	"app/internal/repositories"
)

type EmployeeService struct {
	repo repositories.EmployeeRepository
}

func NewEmployeeService(repo repositories.EmployeeRepository) *EmployeeService {
	return &EmployeeService{
		repo: repo,
	}
}

func (s *EmployeeService) CreateEmployee(ctx context.Context, e *models.Employee) error {
	return s.repo.Create(ctx, e)
}

func (s *EmployeeService) GetByID(ctx context.Context, id int64) (*models.Employee, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *EmployeeService) GetByDepartmentID(ctx context.Context, departmentID int64) ([]*models.Employee, error) {
	return s.repo.FindByDepartmentID(ctx, departmentID)
}
