package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/SeiFlow-3P2/payment_service/internal/models"
	"github.com/SeiFlow-3P2/payment_service/internal/repository"

	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/subscription"
)

type SubscriptionService struct {
	repo repository.SubscriptionRepository
}

func NewSubscriptionService(repo repository.SubscriptionRepository) *SubscriptionService {
	return &SubscriptionService{repo: repo}
}

// Получение локальной записи подписки из БД
func (s *SubscriptionService) GetSubscriptionInfo(ctx context.Context, userID uuid.UUID) (*models.UserSubscription, error) {
	return s.repo.GetByUserID(ctx, userID)
}

// Обновление статуса в локной БД
func (s *SubscriptionService) UpdateStatus(ctx context.Context, userID uuid.UUID, status string) error {
	return s.repo.UpdateStatus(ctx, userID, status)
}

// Получение актуальной информации о подписке из Stripe по userID
func (s *SubscriptionService) GetCurrentSubscription(ctx context.Context, userID uuid.UUID) (*stripe.Subscription, error) {
	subRecord, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if subRecord.StripeSubscriptionID == "" {
		return nil, errors.New("no Stripe subscription ID found")
	}

	subscription, err := subscription.Get(subRecord.StripeSubscriptionID, nil)
	if err != nil {
		return nil, err
	}

	// Проверим, есть ли хотя бы один элемент в подписке
	if len(subscription.Items.Data) == 0 {
		return nil, errors.New("subscription has no items")
	}

	// Пример: получаем имя плана
	planID := subscription.Items.Data[0].Plan.ID
	_ = planID // или используй это значение в ответе

	return subscription, nil
}
