package repositories

import (
	"context"
	"database/sql"
	"strconv"
	"strings"

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
	List(ctx context.Context, limit, offset int, departmentID *int64, keyword *string) ([]*models.Employee, int64, error)
	Update(ctx context.Context, e *models.Employee) error
	Delete(ctx context.Context, id int64) error
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

func (r *employeePostgresRepository) List(ctx context.Context, limit, offset int, departmentID *int64, keyword *string) ([]*models.Employee, int64, error) {
	whereParts := []string{}
	args := []interface{}{}
	if departmentID != nil {
		whereParts = append(whereParts, "department_id = $"+strconv.Itoa(len(args)+1))
		args = append(args, *departmentID)
	}
	if keyword != nil && *keyword != "" {
		whereParts = append(whereParts, "(name ILIKE $"+strconv.Itoa(len(args)+1)+" OR position ILIKE $"+strconv.Itoa(len(args)+1)+")")
		args = append(args, "%"+*keyword+"%")
	}

	where := ""
	if len(whereParts) > 0 {
		where = "WHERE " + strings.Join(whereParts, " AND ")
	}

	var total int64
	countQuery := "SELECT COUNT(*) FROM employees " + where
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	argPos := len(args) + 1
	args = append(args, limit, offset)
	query := "SELECT id, name, email, department_id, age, position, salary, created_at, updated_at FROM employees " + where +
		" ORDER BY id DESC LIMIT $" + strconv.Itoa(argPos) + " OFFSET $" + strconv.Itoa(argPos+1)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var res []*models.Employee
	for rows.Next() {
		var e models.Employee
		var email sql.NullString
		if err := rows.Scan(&e.ID, &e.Name, &email, &e.DepartmentID, &e.Age, &e.Position, &e.Salary, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, 0, err
		}
		if email.Valid {
			e.Email = &email.String
		} else {
			e.Email = nil
		}
		res = append(res, &e)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return res, total, nil
}

func (r *employeePostgresRepository) Update(ctx context.Context, e *models.Employee) error {
	var email sql.NullString
	if e.Email != nil {
		email = sql.NullString{String: *e.Email, Valid: true}
	}

	query := `UPDATE employees SET name = $1, email = $2, department_id = $3, age = $4, position = $5, salary = $6, updated_at = now() WHERE id = $7 RETURNING updated_at`
	var updatedAt sql.NullTime
	if err := r.db.QueryRowContext(ctx, query, e.Name, email, e.DepartmentID, e.Age, e.Position, e.Salary, e.ID).Scan(&updatedAt); err != nil {
		return err
	}
	if updatedAt.Valid {
		e.UpdatedAt = updatedAt.Time
	}
	return nil
}

func (r *employeePostgresRepository) Delete(ctx context.Context, id int64) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM employees WHERE id = $1`, id)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}