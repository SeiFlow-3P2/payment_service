package repository

import (
	"context"
	"payment_service/internal/models"
)

type SubscriptionRepository interface {
	GetByUserID(ctx context.Context, userID string) (*models.UserSubscription, error)
	CreateOrUpdate(ctx context.Context, sub *models.UserSubscription) error
	UpdateStatus(ctx context.Context, userID string, status string) error
}
