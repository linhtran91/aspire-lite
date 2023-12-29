package repositories

import (
	"aspire-lite/internals/constants"
	"aspire-lite/internals/models"
	"context"

	"gorm.io/gorm"
)

type repayment struct {
	db *gorm.DB
}

func NewRepayment(db *gorm.DB) *repayment {
	return &repayment{db: db}
}

func (r *repayment) SubmitRepayment(ctx context.Context, repayment *models.Repayment) error {
	if err := r.db.WithContext(ctx).
		Model(&models.Repayment{}).
		Where("id = ?", repayment.ID).
		Updates(map[string]interface{}{
			"actual_amount": repayment.ActualAmount,
			"paid_at":       repayment.PaidAt,
			"updated_at":    repayment.UpdatedAt,
			"status":        repayment.Status,
		}).Error; err != nil {
		return err
	}
	return nil
}

func (r *repayment) GetByID(ctx context.Context, id string) (*models.Repayment, error) {
	var repayment models.Repayment
	if err := r.db.WithContext(ctx).
		Model(&models.Repayment{}).
		Where("id = ?", id).
		First(&repayment).
		Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, constants.ErrorRecordNotFound
		}
		return nil, err
	}
	return &repayment, nil
}

func (r *repayment) CountUnpaidRepayment(ctx context.Context, loanID int64) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&models.Repayment{}).
		Where("loan_id = ?", loanID).
		Where("status != ?", constants.PAID).
		Count(&count).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return -1, constants.ErrorRecordNotFound
		}
		return -1, err
	}
	return count, nil
}
