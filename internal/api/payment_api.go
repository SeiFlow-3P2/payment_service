package api

import (
	"io"
	"net/http"
	"context"
	"errors"
	"fmt"
	
	"github.com/google/uuid"
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

type WebhookHandler struct {
    paymentService *service.PaymentService
    shutdownChan   chan struct{}
}

func NewWebhookHandler(paymentService *service.PaymentService, shutdownChan chan struct{}) *WebhookHandler {
    return &WebhookHandler{
        paymentService: paymentService,
        shutdownChan:   shutdownChan,
    }
}

func (h *WebhookHandler) HandleStripeWebhook(w http.ResponseWriter, r *http.Request) {
    const MaxBodyBytes = int64(102400) // 100 КБ
    r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)
    defer r.Body.Close()

    // Проверка наличия заголовка Stripe-Signature
    if r.Header.Get("Stripe-Signature") == "" {
        http.Error(w, "Missing Stripe-Signature header", http.StatusBadRequest)
        return
    }

    // Чтение тела запроса
    payload, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "Read error", http.StatusServiceUnavailable)
        return
    }

    // Обработка вебхука через paymentService
    if err := h.paymentService.HandleStripeWebhook(r.Context(), payload, r.Header.Get("Stripe-Signature")); err != nil {
        if err.Error() == "invalid signature" {
            http.Error(w, "Invalid Stripe signature", http.StatusBadRequest)
        } else {
            http.Error(w, "Webhook handling failed: "+err.Error(), http.StatusInternalServerError)
        }
        return
    }

    // Успешный ответ для Stripe
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Webhook processed successfully"))

    // Отправка сигнала о завершении работы вебхука
go func() {
    h.shutdownChan <- struct{}{}
}()

}
// Получение текущей Stripe-подписки пользователя
func (api *PaymentAPI) GetCurrentSubscription(ctx context.Context, req *pb.GetCurrentSubscriptionRequest) (*pb.GetCurrentSubscriptionResponse, error) {
	// Convert string to UUID
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID format: %v", err)
	}
	sub, err := api.subscriptionService.GetCurrentSubscription(ctx, userID)
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
