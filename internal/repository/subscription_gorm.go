package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/SeiFlow-3P2/payment_service/internal/models"
	"gorm.io/gorm"
)

// SubscriptionRepository определяет интерфейс для работы с подписками
type SubscriptionRepository interface {
	// Получить последнюю подписку пользователя
	GetByUserID(ctx context.Context, userID uuid.UUID) (*models.UserSubscription, error)

	// Создать или обновить подписку
	CreateOrUpdate(ctx context.Context, sub *models.UserSubscription) error

	// Обновить статус подписки
	UpdateStatus(ctx context.Context, userID uuid.UUID, status string) error
}

// subscriptionRepositoryGorm — реализация интерфейса SubscriptionRepository с использованием GORM
type subscriptionRepositoryGorm struct {
	db *gorm.DB
}

// NewSubscriptionGorm создает новый экземпляр SubscriptionRepository
func NewSubscriptionGorm(db *gorm.DB) SubscriptionRepository {
	return &subscriptionRepositoryGorm{db: db}
}

// GetByUserID получает подписку по user_id
func (r *subscriptionRepositoryGorm) GetByUserID(ctx context.Context, userID uuid.UUID) (*models.UserSubscription, error) {
	var subscription models.UserSubscription
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID.String()).First(&subscription).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &subscription, nil
}

// CreateOrUpdate создает или обновляет подписку
func (r *subscriptionRepositoryGorm) CreateOrUpdate(ctx context.Context, sub *models.UserSubscription) error {
	return r.db.WithContext(ctx).Save(sub).Error
}

// UpdateStatus обновляет статус подписки
func (r *subscriptionRepositoryGorm) UpdateStatus(ctx context.Context, userID uuid.UUID, status string) error {
	return r.db.WithContext(ctx).Model(&models.UserSubscription{}).
		Where("user_id = ?", userID.String()).
		Update("status", status).Error
}