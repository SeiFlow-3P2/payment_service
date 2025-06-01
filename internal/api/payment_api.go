package api

import (
	"context"
	"errors"

	"github.com/SeiFlow-3P2/payment_service/internal/service"
	pb "github.com/SeiFlow-3P2/payment_service/pkg/proto/v1"
)

type PaymentAPI struct {
	pb.UnimplementedPaymentServiceServer
	paymentService      *service.PaymentService
	subscriptionService *service.SubscriptionService
}

func NewPaymentAPI(paymentService *service.PaymentService, subscriptionService *service.SubscriptionService) *PaymentAPI {
	return &PaymentAPI{
		paymentService:      paymentService,
		subscriptionService: subscriptionService,
	}
}

func (api *PaymentAPI) CreateCheckoutSession(ctx context.Context, req *pb.CreateCheckoutSessionRequest) (*pb.CreateCheckoutSessionResponse, error) {
	session, err := api.paymentService.CreateStripeCheckoutSession(ctx, req.PlanId, req.SuccessUrl, req.CancelUrl)
	if err != nil {
		return nil, err
	}

	return &pb.CreateCheckoutSessionResponse{
		CheckoutSessionId: session.ID,
		CheckoutUrl:       session.URL,
	}, nil
}

// Получение текущей Stripe-подписки пользователя
func (api *PaymentAPI) GetCurrentSubscription(ctx context.Context, req *pb.GetCurrentSubscriptionRequest) (*pb.GetCurrentSubscriptionResponse, error) {
	sub, err := api.subscriptionService.GetCurrentSubscription(ctx, int(req.UserId))
	if err != nil {
		return nil, err
	}
	if sub == nil || len(sub.Items.Data) == 0 {
		return nil, errors.New("subscription not found or has no items")
	}

	item := sub.Items.Data[0]
	if item.Plan == nil {
		return nil, errors.New("plan info is missing in subscription item")
	}

	return &pb.GetCurrentSubscriptionResponse{
		Status:           string(sub.Status),
		PlanId:           item.Plan.ID,
		CurrentPeriodEnd: sub.CurrentPeriodEnd,
	}, nil
}
