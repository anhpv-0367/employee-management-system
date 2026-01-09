package repositories

import (
	"context"
  "database/sql"

	"app/internal/models"
)

type employeePostgresRepository struct {
	db *sql.DB
}

func NewEmployeeRepository(db *sql.DB) EmployeeRepository {
	return &employeePostgresRepository{db: db}
}

type EmployeeRepository interface {
	Create(ctx context.Context, e *models.Employee) error
	FindByID(ctx context.Context, id int64) (*models.Employee, error)
	FindByDepartmentID(ctx context.Context, departmentID int64) ([]*models.Employee, error)
}

func (r *employeePostgresRepository) Create(ctx context.Context, e *models.Employee) error {
	query := `
		INSERT INTO employees (name, email, department_id)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	return r.db.QueryRowContext(
		ctx,
		query,
		e.Name,
		e.Email,
		e.DepartmentID,
	).Scan(&e.ID)
}

func (r *employeePostgresRepository) FindByID(ctx context.Context, id int64) (*models.Employee, error) {
	query := `
		SELECT
			id,
			name,
			email,
			department_id,
			age,
			position,
			salary,
			created_at,
			updated_at
		FROM employees
		WHERE id = $1
	`

	var e models.Employee
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&e.ID,
		&e.Name,
		&e.Email,
		&e.DepartmentID,
		&e.Age,
		&e.Position,
		&e.Salary,
		&e.CreatedAt,
		&e.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &e, nil
}

func (r *employeePostgresRepository) FindByDepartmentID(ctx context.Context, departmentID int64) ([]*models.Employee, error) {
	query := `
		SELECT
			id,
			name,
			email,
			department_id,
			age,
			position,
			salary,
			created_at,
			updated_at
		FROM employees
		WHERE department_id = $1
	`

	rows, err := r.db.QueryContext(ctx, query, departmentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var employees []*models.Employee
	for rows.Next() {
		var e models.Employee
		if err := rows.Scan(&e.ID, &e.Name, &e.Email, &e.DepartmentID); err != nil {
			return nil, err
		}
		employees = append(employees, &e)
	}
	return employees, nil
}