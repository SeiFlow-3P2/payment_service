package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/SeiFlow-3P2/payment_service/internal/api"
	"github.com/SeiFlow-3P2/payment_service/internal/models"
	"github.com/SeiFlow-3P2/payment_service/internal/repository"
	"github.com/SeiFlow-3P2/payment_service/internal/service"
	pb "github.com/SeiFlow-3P2/payment_service/pkg/proto/v1"
	"google.golang.org/grpc"
)

func main() {
	// Load .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file:", err)
	}

	// Check required environment variables
	if os.Getenv("STRIPE_SECRET_KEY") == "" {
		log.Fatal("STRIPE_SECRET_KEY not found in environment variables")
	}
	if os.Getenv("STRIPE_WEBHOOK_SECRET") == "" {
		log.Fatal("STRIPE_WEBHOOK_SECRET not found in environment variables")
	}

	// Initialize database connection
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		// Переменная окружения не задана — используем локальную БД
		dsn = "postgres://postgres:12345@localhost:5432/payment?sslmode=disable"
		log.Println("DATABASE_URL not found, using local database:", dsn)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto-migrate the schema
	if err := db.AutoMigrate(&models.PaymentRecord{}, &models.UserSubscription{}); err != nil {
		log.Fatal("Failed to migrate database schema:", err)
	}

	// Initialize dependencies
	repo := repository.NewPaymentRecordGorm(db)
	paymentService := service.NewPaymentService(repo)
	paymentAPI := api.NewPaymentAPI(paymentService)

	// Start gRPC server
	go func() {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		grpcServer := grpc.NewServer()
		pb.RegisterPaymentServiceServer(grpcServer, paymentAPI)
		log.Println("gRPC server started on :50051")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Start HTTP server with Stripe Webhook endpoint
	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		// Логируем заголовки, чтобы проверить наличие "Stripe-Signature"
		for name, values := range r.Header {
			log.Printf("Header: %s=%s", name, values)
		}

		payload, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// Получаем Stripe-Signature из заголовков
		sig := r.Header.Get("Stripe-Signature")
		fmt.Println(sig)
		if sig == "" {
			log.Println("Stripe-Signature header missing")
			http.Error(w, "Missing Stripe-Signature header", http.StatusBadRequest)
			return
		}

		// Обрабатываем событие с проверкой подписи
		err = paymentService.HandleStripeWebhook(context.Background(), payload, sig)
		if err != nil {
			log.Println("Webhook error:", err)
			http.Error(w, "Webhook Error: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Ответ на успешную обработку
		w.WriteHeader(http.StatusOK)
	})

	log.Println("HTTP server started on :80 (webhook endpoint)")
	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Fatalf("HTTP server error: %v", err)
	}
}
