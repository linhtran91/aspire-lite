package models

import "time"

type Loan struct {
	ID            int64
	Amount        float64
	Term          int
	CustomerID    int64
	Status        int8
	ScheduledDate time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
