package repository

import (
	"context"

	"github.com/SeiFlow-3P2/payment_service/internal/models"
)

type SubscriptionRepository interface {
	// Получить последнюю подписку пользователя
	GetByUserID(ctx context.Context, userID int) (*models.UserSubscription, error)

	// Создать или обновить подписку
	CreateOrUpdate(ctx context.Context, sub *models.UserSubscription) error

	// Обновить статус подписки
	UpdateStatus(ctx context.Context, userID string, status string) error
}
