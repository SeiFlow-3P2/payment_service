package repository

import (
	"context"
	"payment_service/internal/models"
)

type PaymentRecordRepository interface {
	Create(ctx context.Context, record *models.PaymentRecord) error
	UpdateStatus(ctx context.Context, checkoutSessionID string, status string, chargeID string) error
	GetByCheckoutSessionID(ctx context.Context, checkoutSessionID string) (*models.PaymentRecord, error)
}
