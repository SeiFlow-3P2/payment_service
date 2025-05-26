package service

import (
	"context"
	"encoding/json"
	"errors"
	"os"

	"github.com/SeiFlow-3P2/payment_service/internal/models"
	"github.com/SeiFlow-3P2/payment_service/internal/repository"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/checkout/session"
	"github.com/stripe/stripe-go/v76/webhook"
)

type PaymentService struct {
	repo repository.PaymentRecordRepository
}

func NewPaymentService(repo repository.PaymentRecordRepository) *PaymentService {
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
	return &PaymentService{repo: repo}
}

// Бизнес-логика создания Stripe Checkout Session
func (s *PaymentService) CreateStripeCheckoutSession(ctx context.Context, planID, successURL, cancelURL string) (*stripe.CheckoutSession, error) {
	params := &stripe.CheckoutSessionParams{
		SuccessURL: stripe.String(successURL),
		CancelURL:  stripe.String(cancelURL),
		Mode:       stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(planID),
				Quantity: stripe.Int64(1),
			},
		},
	}

	sess, err := session.New(params)
	if err != nil {
		return nil, err
	}

	record := &models.PaymentRecord{
		StripeCheckoutSessionID: sess.ID,
		Status:                  "created",
		PlanID:                  planID,
	}

	err = s.repo.Create(ctx, record)
	if err != nil {
		return nil, err
	}

	return sess, nil
}

// Обновление статуса платежа (используется вебхуком)
func (s *PaymentService) HandleCheckoutCompleted(ctx context.Context, session *stripe.CheckoutSession) error {
	if session.PaymentIntent == nil {
		return errors.New("missing payment intent")
	}
	return s.repo.UpdateStatus(ctx, session.ID, "paid", session.PaymentIntent.ID)
}

// Бизнес-логика обработки webhook: парсинг и делегирование
func (s *PaymentService) HandleStripeWebhook(ctx context.Context, payload []byte, sigHeader string) error {
	endpointSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
	if endpointSecret == "" {
		return errors.New("missing STRIPE_WEBHOOK_SECRET")
	}

	event, err := webhook.ConstructEvent(payload, sigHeader, endpointSecret)
	if err != nil {
		return err
	}

	switch event.Type {
	case "checkout.session.completed":
		var session stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
			return err
		}
		return s.HandleCheckoutCompleted(ctx, &session)
	default:
		// Необработанный тип события не является ошибкой
		return nil
	}
}

// CRUD методы
func (s *PaymentService) CreatePaymentRecord(ctx context.Context, record *models.PaymentRecord) error {
	return s.repo.Create(ctx, record)
}

func (s *PaymentService) UpdatePaymentStatus(ctx context.Context, checkoutSessionID, status, chargeID string) error {
	return s.repo.UpdateStatus(ctx, checkoutSessionID, status, chargeID)
}

func (s *PaymentService) GetPaymentRecord(ctx context.Context, checkoutSessionID string) (*models.PaymentRecord, error) {
	return s.repo.GetByCheckoutSessionID(ctx, checkoutSessionID)
}
