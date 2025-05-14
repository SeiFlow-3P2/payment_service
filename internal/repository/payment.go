package repository

import (
	"context"
	"payment_service/internal/models"
)

type PaymentRepository interface {
	GetByID(ctx context.Context, id uint) (*models.Payment, error)
	Create(ctx context.Context, payment *models.Payment) error
	Update(ctx context.Context, payment *models.Payment) error
	Delete(ctx context.Context, id uint) error
}
