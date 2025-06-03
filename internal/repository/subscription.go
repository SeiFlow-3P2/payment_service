package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/SeiFlow-3P2/payment_service/internal/models"
)

type SubscriptionRepository interface {
	// Получить последнюю подписку пользователя
	GetByUserID(ctx context.Context, userID uuid.UUID) (*models.UserSubscription, error)

	// Создать или обновить подписку
	CreateOrUpdate(ctx context.Context, sub *models.UserSubscription) error

	// Обновить статус подписки
	UpdateStatus(ctx context.Context, userID uuid.UUID, status string) error
}
