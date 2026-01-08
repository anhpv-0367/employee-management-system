package models

import "time"

type Employee struct {
	ID           int64
	Name         string
	Email        *string
	DepartmentID int64
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
