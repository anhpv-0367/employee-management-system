package services

import (
	"context"

	"app/internal/models"
	"app/internal/repositories"
)

type DepartmentService struct {
	repo repositories.DepartmentRepository
}

func NewDepartmentService(repo repositories.DepartmentRepository) *DepartmentService {
	return &DepartmentService{
		repo: repo,
	}
}

func (s *DepartmentService) Create(ctx context.Context, d *models.Department) error {
	return s.repo.Create(ctx, d)
}
