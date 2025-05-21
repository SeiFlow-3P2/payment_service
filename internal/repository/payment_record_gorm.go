package repository

import (
	"context"

	"github.com/SeiFlow-3P2/payment_service/internal/models"

	"gorm.io/gorm"
)

type paymentRecordGorm struct {
	db *gorm.DB
}

func NewPaymentRecordGorm(db *gorm.DB) PaymentRecordRepository {
	return &paymentRecordGorm{db: db}
}

func (r *paymentRecordGorm) Create(ctx context.Context, record *models.PaymentRecord) error {
	return r.db.WithContext(ctx).Create(record).Error
}

func (r *paymentRecordGorm) UpdateStatus(ctx context.Context, checkoutSessionID string, status string, chargeID string) error {
	return r.db.WithContext(ctx).
		Model(&models.PaymentRecord{}).
		Where("stripe_checkout_session_id = ?", checkoutSessionID).
		Updates(map[string]interface{}{
			"status":           status,
			"stripe_charge_id": chargeID,
		}).Error
}

func (r *paymentRecordGorm) GetByCheckoutSessionID(ctx context.Context, checkoutSessionID string) (*models.PaymentRecord, error) {
	var record models.PaymentRecord
	if err := r.db.WithContext(ctx).
		Where("stripe_checkout_session_id = ?", checkoutSessionID).
		First(&record).Error; err != nil {
		return nil, err
	}
	return &record, nil
}
