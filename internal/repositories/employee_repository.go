package repositories

import (
	"context"

	"app/internal/models"
)

type EmployeeRepository interface {
	Create(ctx context.Context, e *models.Employee) error
	FindByID(ctx context.Context, id int64) (*models.Employee, error)
	FindByDepartmentID(ctx context.Context, departmentID int64) ([]*models.Employee, error)
}
