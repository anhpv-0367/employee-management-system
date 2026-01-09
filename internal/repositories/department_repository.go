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
	FindAll(ctx context.Context) ([]*models.Department, error)
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

func (r *departmentPostgresRepository) FindAll(ctx context.Context) ([]*models.Department, error) {
	query := `
		SELECT id, name
		FROM departments
		ORDER BY id
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var departments []*models.Department
	for rows.Next() {
		var d models.Department
		if err := rows.Scan(&d.ID, &d.Name); err != nil {
			return nil, err
		}
		departments = append(departments, &d)
	}

	return departments, nil
}
