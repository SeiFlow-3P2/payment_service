package service

import (
	"context"
	"payment_service/internal/models"
	"payment_service/internal/repository"
)

type PaymentService struct {
	repo repository.PaymentRepository
}

func NewPaymentService(repo repository.PaymentRepository) *PaymentService {
	return &PaymentService{repo: repo}
}

func (s *PaymentService) GetPaymentInfo(ctx context.Context, id uint) (*models.Payment, error) {
	return s.repo.GetByID(ctx, id)
}
