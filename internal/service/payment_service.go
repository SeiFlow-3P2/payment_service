package service

import (
	"context"
	"payment_service/internal/models"
	"payment_service/internal/repository"
)

type PaymentService struct {
	repo repository.PaymentRecordRepository
}

func NewPaymentService(repo repository.PaymentRecordRepository) *PaymentService {
	return &PaymentService{repo: repo}
}

func (s *PaymentService) CreatePaymentRecord(ctx context.Context, record *models.PaymentRecord) error {
	return s.repo.Create(ctx, record)
}

func (s *PaymentService) UpdatePaymentStatus(ctx context.Context, checkoutSessionID, status, chargeID string) error {
	return s.repo.UpdateStatus(ctx, checkoutSessionID, status, chargeID)
}

func (s *PaymentService) GetPaymentRecord(ctx context.Context, checkoutSessionID string) (*models.PaymentRecord, error) {
	return s.repo.GetByCheckoutSessionID(ctx, checkoutSessionID)
}
