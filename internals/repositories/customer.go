package repositories

import (
	"aspire-lite/internals/constants"
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
	if err := r.db.WithContext(ctx).Model(&models.Customer{}).Where("username = ?", username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, constants.ErrorRecordNotFound
		}
		return nil, err
	}
	return user, nil
}

func (r *customer) GetUserByID(ctx context.Context, id int64) (*models.Customer, error) {
	var user *models.Customer
	if err := r.db.WithContext(ctx).Model(&models.Customer{}).Where("id = ?", id).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, constants.ErrorRecordNotFound
		}
		return nil, err
	}
	return user, nil
}
