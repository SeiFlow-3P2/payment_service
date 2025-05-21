package api

import (
	"context"
	"os"

	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/checkout/session"

	"github.com/SeiFlow-3P2/payment_service/internal/models"
	"github.com/SeiFlow-3P2/payment_service/internal/service"
	pb "github.com/SeiFlow-3P2/payment_service/pkg/proto/v1"
)

type PaymentAPI struct {
	pb.UnimplementedPaymentServiceServer
	paymentService *service.PaymentService
}

func NewPaymentAPI(paymentService *service.PaymentService) *PaymentAPI {
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
	return &PaymentAPI{paymentService: paymentService}
}

func (api *PaymentAPI) CreateCheckoutSession(ctx context.Context, req *pb.CreateCheckoutSessionRequest) (*pb.CreateCheckoutSessionResponse, error) {
	params := &stripe.CheckoutSessionParams{
		SuccessURL: stripe.String(req.SuccessUrl),
		CancelURL:  stripe.String(req.CancelUrl),
		Mode:       stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(req.PlanId), // price_xxx ID from Stripe
				Quantity: stripe.Int64(1),
			},
		},
	}

	sess, err := session.New(params)
	if err != nil {
		return nil, err
	}

	// Запись в базу данных
	record := &models.PaymentRecord{
		StripeCheckoutSessionID: sess.ID,
		Status:            "created",
		PlanID:            req.PlanId,
	}
	if err := api.paymentService.CreatePaymentRecord(ctx, record); err != nil {
		return nil, err
	}

	return &pb.CreateCheckoutSessionResponse{
		CheckoutSessionId: sess.ID,
		CheckoutUrl:       sess.URL,
	}, nil
}
