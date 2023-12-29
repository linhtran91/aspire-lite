package repositories

import (
	"aspire-lite/internals/models"
	"context"

	"gorm.io/gorm"
)

type customer struct {
	db *gorm.DB
}

func NewCustomer(db *gorm.DB) *customer {
	return &customer{db: db}
}

func (r *customer) GetUserCredential(ctx context.Context, username string) (*models.Customer, error) {
	var user *models.Customer
	if err := r.db.WithContext(ctx).Model(&models.Customer{}).First(&user).Where("username = ?", username).Error; err != nil {
		return nil, err
	}
	return user, nil
}
