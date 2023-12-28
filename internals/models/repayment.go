package models

import "time"

type Repayment struct {
	ID              string
	LoanID          int64
	ScheduledAmount float64
	ActualAmount    float64
	Status          int
	PaidAt          time.Time
	ScheduledPayAt  time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
