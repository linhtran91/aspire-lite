package models

import "time"

type Repayment struct {
	ID              string    `json:"id"`
	LoanID          int64     `json:"loan_id"`
	ScheduledAmount float64   `json:"schedule_date"`
	ActualAmount    float64   `json:"actual_amount"`
	Status          int       `json:"status"`
	PaidAt          time.Time `json:"paid_at"`
	ScheduledPayAt  time.Time `json:"schedule_pay_at"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
