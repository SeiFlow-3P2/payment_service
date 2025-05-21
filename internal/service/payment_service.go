package service

import (
	"context"

	"github.com/SeiFlow-3P2/payment_service/internal/models"
	"github.com/SeiFlow-3P2/payment_service/internal/repository"
	"github.com/stripe/stripe-go/v76"
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

func (s *PaymentService) HandleCheckoutCompleted(ctx context.Context, session *stripe.CheckoutSession) error {
	return s.repo.UpdateStatus(ctx, session.ID, "paid", session.PaymentIntent.ID)
}
