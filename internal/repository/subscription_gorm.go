package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/SeiFlow-3P2/payment_service/internal/models"
	"gorm.io/gorm"
)

type subscriptionGorm struct {
	db *gorm.DB
}

func NewSubscriptionGorm(db *gorm.DB) SubscriptionRepository {
	return &subscriptionGorm{db: db}
}


// Получение последней по времени подписки пользователя
func (r *subscriptionGorm) GetByUserID(ctx context.Context, userID uuid.UUID) (*models.UserSubscription, error) {
	var sub models.UserSubscription
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&sub).Error; err != nil {
		return nil, err
	}
	return &sub, nil
}

func (r *subscriptionGorm) CreateOrUpdate(ctx context.Context, sub *models.UserSubscription) error {
	return r.db.WithContext(ctx).Save(sub).Error
}

func (r *subscriptionGorm) UpdateStatus(ctx context.Context, userID uuid.UUID, status string) error {
	return r.db.WithContext(ctx).
		Model(&models.UserSubscription{}).
		Where("user_id = ?", userID).
		Update("status", status).Error
}

