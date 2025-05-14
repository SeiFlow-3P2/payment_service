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

func (s *SubscriptionService) GetSubscriptionInfo(ctx context.Context, userID string) (*models.UserSubscription, error) {
	return s.repo.GetByUserID(ctx, userID)
}

func (s *SubscriptionService) UpdateStatus(ctx context.Context, userID, status string) error {
	return s.repo.UpdateStatus(ctx, userID, status)
}
