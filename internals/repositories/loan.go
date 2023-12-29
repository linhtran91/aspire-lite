package repositories

import (
	"aspire-lite/internals/constants"
	"aspire-lite/internals/models"
	"context"
	"time"

	"gorm.io/gorm"
)

type loan struct {
	db *gorm.DB
}

func NewLoan(db *gorm.DB) *loan {
	return &loan{db: db}
}

func (r *loan) Create(ctx context.Context, loan *models.Loan, repayments []*models.Repayment) (int64, error) {
	if err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&loan).Error; err != nil {
			return err
		}

		for i := range repayments {
			repayments[i].LoanID = loan.ID
		}

		if err := tx.CreateInBatches(repayments, 50).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return -1, err
	}
	return loan.ID, nil
}

func (r *loan) Approve(ctx context.Context, loanID int64, at time.Time) error {
	if err := r.db.WithContext(ctx).
		Model(&models.Loan{}).
		Where("id = ?", loanID).
		Updates(map[string]interface{}{
			"status":     constants.APPROVED,
			"updated_at": at,
		}).Error; err != nil {
		return err
	}
	return nil
}

func (r *loan) View(ctx context.Context, customerID int64, limit, offset int) ([]*models.Loan, error) {
	var result []*models.Loan
	err := r.db.WithContext(ctx).
		Model(&models.Loan{}).
		Find(&result).
		Where("customer_id = ?", customerID).
		Order("schedule_date desc").
		Limit(limit).Offset(offset).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *loan) UpdateStatus(ctx context.Context, loanID int64) error {
	if err := r.db.WithContext(ctx).Model(&models.Loan{}).
		Where("id = ?", loanID).
		Updates(map[string]interface{}{
			"status": constants.PAID,
		}).Error; err != nil {
		return err
	}
	return nil
}
