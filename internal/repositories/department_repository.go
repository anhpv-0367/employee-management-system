package repositories

import (
	"context"
	"database/sql"
  "errors"

	"app/internal/models"
)

type departmentPostgresRepository struct {
	db *sql.DB
}

func NewDepartmentRepository(db *sql.DB) DepartmentRepository {
	return &departmentPostgresRepository{db: db}
}

type DepartmentRepository interface {
	Create(ctx context.Context, d *models.Department) error
	FindByID(ctx context.Context, id int64) (*models.Department, error)
	FindAll(ctx context.Context, limit, offset int) ([]*models.Department, int64, error)
}

func (r *departmentPostgresRepository) Create(ctx context.Context, d *models.Department) error {
	query := `
		INSERT INTO departments (name)
		VALUES ($1)
		RETURNING id
	`
	return r.db.QueryRowContext(ctx, query, d.Name).Scan(&d.ID)
}

func (r *departmentPostgresRepository) FindByID(ctx context.Context, id int64) (*models.Department, error) {
	query := `
		SELECT id, name
		FROM departments
		WHERE id = $1
	`

	var d models.Department
	err := r.db.QueryRowContext(ctx, query, id).Scan(&d.ID, &d.Name)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	return &d, err
}

func (r *departmentPostgresRepository) FindAll(ctx context.Context, limit, offset int) ([]*models.Department, int64, error) {
	var total int64
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM departments`).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `SELECT id, name, created_at, updated_at FROM departments ORDER BY id LIMIT $1 OFFSET $2`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var departments []*models.Department
	for rows.Next() {
		var d models.Department
		if err := rows.Scan(&d.ID, &d.Name, &d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, 0, err
		}
		departments = append(departments, &d)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return departments, total, nil
}
