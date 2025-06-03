package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/SeiFlow-3P2/payment_service/internal/models"
	"github.com/SeiFlow-3P2/payment_service/internal/repository"
)

type SubscriptionService struct {
	repo repository.SubscriptionRepository // Теперь это интерфейс
}

func NewSubscriptionService(repo repository.SubscriptionRepository) *SubscriptionService {
	return &SubscriptionService{repo: repo}
}

func (s *SubscriptionService) GetSubscriptionByUserID(ctx context.Context, userID uuid.UUID) (*models.UserSubscription, error) {
	if userID == uuid.Nil {
		return nil, errors.New("invalid user ID")
	}

	// Используем метод интерфейса GetByUserID
	subscription, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if subscription == nil {
		return nil, errors.New("subscription not found")
	}

	return subscription, nil
}