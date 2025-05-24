package api

import (
	"context"

	"github.com/SeiFlow-3P2/payment_service/internal/service"
	pb "github.com/SeiFlow-3P2/payment_service/pkg/proto/v1"
)

type PaymentAPI struct {
	pb.UnimplementedPaymentServiceServer
	paymentService *service.PaymentService
}

func NewPaymentAPI(paymentService *service.PaymentService) *PaymentAPI {
	return &PaymentAPI{paymentService: paymentService}
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
