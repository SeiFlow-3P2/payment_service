package api

import (
	"io"
	"net/http"
	"os"
	"encoding/json"

	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/webhook"

	"github.com/SeiFlow-3P2/payment_service/internal/service"
)

type WebhookHandler struct {
	paymentService *service.PaymentService
}

func NewWebhookHandler(paymentService *service.PaymentService) *WebhookHandler {
	return &WebhookHandler{paymentService: paymentService}
}

func (h *WebhookHandler) HandleStripeWebhook(w http.ResponseWriter, r *http.Request) {
	const MaxBodyBytes = int64(65536)
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)

	payload, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Read error", http.StatusServiceUnavailable)
		return
	}

	endpointSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
	sigHeader := r.Header.Get("Stripe-Signature")

	event, err := webhook.ConstructEvent(payload, sigHeader, endpointSecret)
	if err != nil {
		http.Error(w, "Signature verification failed", http.StatusBadRequest)
		return
	}

	switch event.Type {
	case "checkout.session.completed":
		var session stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
			http.Error(w, "Failed to parse session", http.StatusBadRequest)
			return
		}

		err := h.paymentService.HandleCheckoutCompleted(r.Context(), &session)
		if err != nil {
			http.Error(w, "Service error", http.StatusInternalServerError)
			return
		}
	default:
	}

	w.WriteHeader(http.StatusOK)
}
