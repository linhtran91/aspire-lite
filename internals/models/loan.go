package models

import "time"

type Loan struct {
	ID            int64     `json:"id"`
	Amount        float64   `json:"amount"`
	Term          int       `json:"term"`
	CustomerID    int64     `json:"customer_id"`
	Status        int8      `json:"status"`
	ScheduledDate time.Time `json:"schedule_date"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
