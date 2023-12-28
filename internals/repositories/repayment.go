package repositories

import (
	"aspire-lite/internals/models"

	"gorm.io/gorm"
)

type repayment struct {
	db *gorm.DB
}

func NewRepayment(db *gorm.DB) *repayment {
	return &repayment{db: db}
}

func (r *repayment) PayByCustomer(repayment *models.Repayment) error {
	return nil
}
