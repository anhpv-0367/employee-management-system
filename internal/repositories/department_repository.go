package repositories

import (
	"context"

	"app/internal/models"
)

type DepartmentRepository interface {
	Create(ctx context.Context, d *models.Department) error
	FindByID(ctx context.Context, id int64) (*models.Department, error)
	FindAll(ctx context.Context) ([]*models.Department, error)
}
