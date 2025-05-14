package service

import (
	"context"
	"payment_service/internal/models"
	"payment_service/internal/repository"
)

type SubscriptionService struct {
	repo repository.SubscriptionRepository
}

func NewSubscriptionService(repo repository.SubscriptionRepository) *SubscriptionService {
	return &SubscriptionService{repo: repo}
}

func (s *SubscriptionService) GetSubscriptionInfo(ctx context.Context, id uint) (*models.Subscription, error) {
	return s.repo.GetByID(ctx, id)
}
