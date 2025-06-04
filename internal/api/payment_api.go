package api

import (

	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"github.com/SeiFlow-3P2/payment_service/internal/service"
	pb "github.com/SeiFlow-3P2/payment_service/pkg/proto/v1"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
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
func (api *PaymentAPI) GetSubscriptionInfo(ctx context.Context, req *pb.GetSubscriptionInfoRequest) (*pb.GetSubscriptionInfoResponse, error) {
	log.Printf("GetSubscriptionInfo called for user_id: %s", req.UserId)


	// Преобразуем user_id (строка) в uuid.UUID
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		log.Printf("GetSubscriptionInfo error: invalid user ID format: %v", err)
		return nil, errors.New("invalid user ID format")
	}

	// Получаем подписку из базы данных по user_id
	sub, err := api.subscriptionService.GetSubscriptionByUserID(ctx, userID)
	if err != nil {
		log.Printf("GetSubscriptionInfo error: failed to get subscription: %v", err)
		return nil, err
	}
	if sub == nil {
		return nil, errors.New("subscription not found")
	}

	// Преобразуем данные в формат ответа
	response := &pb.GetSubscriptionInfoResponse{
		PlanId:             sub.PlanID,
		Status:             sub.Status,
		CurrentPeriodStart: timestamppb.New(sub.CurrentPeriodStart),
		CurrentPeriodEnd:   timestamppb.New(sub.CurrentPeriodEnd),
	}

	return response, nil
}
