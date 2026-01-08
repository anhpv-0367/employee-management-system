package models

import "time"

type Department struct {
	ID        int64
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
