package services

import (
	"context"
	"errors"

	"app/internal/models"
	"app/internal/repositories"
)

type EmployeeService struct {
	repo     repositories.EmployeeRepository
	deptRepo repositories.DepartmentRepository
}

func NewEmployeeService(repo repositories.EmployeeRepository, deptRepo repositories.DepartmentRepository) *EmployeeService {
	return &EmployeeService{
		repo:     repo,
		deptRepo: deptRepo,
	}
}

func (s *EmployeeService) GetByID(ctx context.Context, id int64) (*models.Employee, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *EmployeeService) GetByDepartmentID(ctx context.Context, departmentID int64) ([]*models.Employee, error) {
	return s.repo.FindByDepartmentID(ctx, departmentID)
}

func (s *EmployeeService) List(ctx context.Context, limit, offset int, departmentID *int64, keyword *string) ([]*models.Employee, int64, error) {
	return s.repo.List(ctx, limit, offset, departmentID, keyword)
}

func (s *EmployeeService) CreateEmployee(ctx context.Context, e *models.Employee) error {
	if e.DepartmentID == 0 {
		return errors.New("departmentId is required")
	}
	if _, err := s.deptRepo.FindByID(ctx, e.DepartmentID); err != nil {
		return errors.New("department not found")
	}
	return s.repo.Create(ctx, e)
}

func (s *EmployeeService) Update(ctx context.Context, e *models.Employee) error {
	if e.DepartmentID == 0 {
		return errors.New("departmentId is required")
	}
	if _, err := s.deptRepo.FindByID(ctx, e.DepartmentID); err != nil {
		return errors.New("department not found")
	}
	return s.repo.Update(ctx, e)
}

func (s *EmployeeService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
