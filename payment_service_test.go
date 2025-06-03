package payment_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	paymentpb "payment_service/pkg/proto/payment/v1"
)

// mockPaymentService implements paymentpb.PaymentServiceServer for testing
type mockPaymentService struct {
	paymentpb.UnimplementedPaymentServiceServer
}

func (m *mockPaymentService) CreateCheckoutSession(ctx context.Context, req *paymentpb.CreateCheckoutSessionRequest) (*paymentpb.CreateCheckoutSessionResponse, error) {
	if req.GetPlanId() == "" || req.GetSuccessUrl() == "" || req.GetCancelUrl() == "" {
		return nil, errors.New("missing required fields")
	}

	return &paymentpb.CreateCheckoutSessionResponse{
		CheckoutSessionId: "mock_session_123",
		CheckoutUrl:       "https://checkout.mock/session_123",
	}, nil
}

func (m *mockPaymentService) HandleStripeWebhook(ctx context.Context, req *paymentpb.HandleStripeWebhookRequest) (*paymentpb.HandleStripeWebhookResponse, error) {
	if req.GetPayload() == "" || req.GetStripeSignature() == "" {
		return nil, errors.New("missing required fields")
	}

	return &paymentpb.HandleStripeWebhookResponse{
		Success: true,
		Message: "webhook processed",
	}, nil
}

func (m *mockPaymentService) GetSubscriptionInfo(ctx context.Context, req *paymentpb.GetSubscriptionInfoRequest) (*paymentpb.GetSubscriptionInfoResponse, error) {
	now := time.Now()
	return &paymentpb.GetSubscriptionInfoResponse{
		PlanId:             "premium_monthly",
		Status:            "active",
		CurrentPeriodStart: timestamppb.New(now),
		CurrentPeriodEnd:   timestamppb.New(now.Add(30 * 24 * time.Hour)),
	}, nil
}

func TestPaymentService(t *testing.T) {
	ctx := context.Background()
	service := &mockPaymentService{}

	t.Run("CreateCheckoutSession", func(t *testing.T) {
		tests := []struct {
			name    string
			req     *paymentpb.CreateCheckoutSessionRequest
			wantErr bool
		}{
			{
				name: "success",
				req: &paymentpb.CreateCheckoutSessionRequest{
					PlanId:     "premium_monthly",
					SuccessUrl: "https://example.com/success",
					CancelUrl:  "https://example.com/cancel",
				},
				wantErr: false,
			},
			{
				name: "missing plan_id",
				req: &paymentpb.CreateCheckoutSessionRequest{
					SuccessUrl: "https://example.com/success",
					CancelUrl:  "https://example.com/cancel",
				},
				wantErr: true,
			},
			{
				name: "missing success_url",
				req: &paymentpb.CreateCheckoutSessionRequest{
					PlanId:    "premium_monthly",
					CancelUrl: "https://example.com/cancel",
				},
				wantErr: true,
			},
			{
				name: "missing cancel_url",
				req: &paymentpb.CreateCheckoutSessionRequest{
					PlanId:     "premium_monthly",
					SuccessUrl: "https://example.com/success",
				},
				wantErr: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				resp, err := service.CreateCheckoutSession(ctx, tt.req)
				if (err != nil) != tt.wantErr {
					t.Errorf("CreateCheckoutSession() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !tt.wantErr && resp.GetCheckoutSessionId() == "" {
					t.Error("expected non-empty checkout session ID")
				}
			})
		}
	})

	t.Run("HandleStripeWebhook", func(t *testing.T) {
		tests := []struct {
			name    string
			req     *paymentpb.HandleStripeWebhookRequest
			wantErr bool
		}{
			{
				name: "success",
				req: &paymentpb.HandleStripeWebhookRequest{
					Payload:        "test_payload",
					StripeSignature: "test_sig",
				},
				wantErr: false,
			},
			{
				name: "missing payload",
				req: &paymentpb.HandleStripeWebhookRequest{
					StripeSignature: "test_sig",
				},
				wantErr: true,
			},
			{
				name: "missing signature",
				req: &paymentpb.HandleStripeWebhookRequest{
					Payload: "test_payload",
				},
				wantErr: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				resp, err := service.HandleStripeWebhook(ctx, tt.req)
				if (err != nil) != tt.wantErr {
					t.Errorf("HandleStripeWebhook() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !tt.wantErr && !resp.GetSuccess() {
					t.Error("expected success response")
				}
			})
		}
	})

	t.Run("GetSubscriptionInfo", func(t *testing.T) {
		resp, err := service.GetSubscriptionInfo(ctx, &paymentpb.GetSubscriptionInfoRequest{})
		if err != nil {
			t.Fatalf("GetSubscriptionInfo() error = %v", err)
		}

		if resp.GetPlanId() != "premium_monthly" {
			t.Errorf("expected plan_id 'premium_monthly', got %q", resp.GetPlanId())
		}

		if resp.GetStatus() != "active" {
			t.Errorf("expected status 'active', got %q", resp.GetStatus())
		}

		if resp.GetCurrentPeriodStart() == nil || resp.GetCurrentPeriodEnd() == nil {
			t.Error("expected non-nil period dates")
		}

		if resp.GetCurrentPeriodEnd().AsTime().Before(resp.GetCurrentPeriodStart().AsTime()) {
			t.Error("expected period_end to be after period_start")
		}
	})
}