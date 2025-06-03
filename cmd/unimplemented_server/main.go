package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/SeiFlow-3P2/payment_service/internal/api"
	"github.com/SeiFlow-3P2/payment_service/internal/models"
	"github.com/SeiFlow-3P2/payment_service/internal/repository"
	"github.com/SeiFlow-3P2/payment_service/internal/service"
	pb "github.com/SeiFlow-3P2/payment_service/pkg/proto/v1"
)

func main() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file:", err)
	}

	if os.Getenv("STRIPE_SECRET_KEY") == "" {
		log.Fatal("STRIPE_SECRET_KEY not found in environment variables")
	}
	if os.Getenv("STRIPE_WEBHOOK_SECRET") == "" {
		log.Fatal("STRIPE_WEBHOOK_SECRET not found in environment variables")
	}

	// Initialize DB
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:12345@localhost:5432/payment?sslmode=disable"
		log.Println("DATABASE_URL not found, using fallback:", dsn)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}

	if err := db.AutoMigrate(&models.PaymentRecord{}, &models.UserSubscription{}); err != nil {
		log.Fatal("Failed DB migration:", err)
	}

	// Dependencies
	repo := repository.NewPaymentRecordGorm(db)
	subscriptionRepo := repository.NewSubscriptionGorm(db)
	paymentService := service.NewPaymentService(repo)
	subscriptionService := service.NewSubscriptionService(subscriptionRepo)
	paymentAPI := api.NewPaymentAPI(paymentService, subscriptionService)

	// gRPC server
	grpcServer := grpc.NewServer()
	pb.RegisterPaymentServiceServer(grpcServer, paymentAPI)

	grpcLis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen on gRPC: %v", err)
	}

	go func() {
		log.Println("gRPC server started on :50052")
		if err := grpcServer.Serve(grpcLis); err != nil {
			log.Fatalf("gRPC server failed: %v", err)
		}
	}()

	// REST Gateway (grpc-gateway) — подключаем маршруты
	go func() {
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		mux := runtime.NewServeMux()
		opts := []grpc.DialOption{grpc.WithInsecure()} // ⚠️ Без TLS

		err := pb.RegisterPaymentServiceHandlerFromEndpoint(ctx, mux, "localhost:50052", opts)
		if err != nil {
			log.Fatalf("Failed to register gRPC-Gateway: %v", err)
		}

		log.Println("REST gateway started on :8080")
		if err := http.ListenAndServe(":8080", mux); err != nil && err != http.ErrServerClosed {
			log.Fatalf("REST gateway failed: %v", err)
		}
	}()

	// Webhook endpoint — отдельно на :80
	webhookHandler := api.NewWebhookHandler(paymentService, make(chan struct{}))

	go func() {
		http.HandleFunc("/webhook", webhookHandler.HandleStripeWebhook)
		log.Println("HTTP webhook server started on :80 (/webhook)")
		if err := http.ListenAndServe(":80", nil); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Webhook server error: %v", err)
		}
	}()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("Shutting down...")

	grpcServer.GracefulStop()
	log.Println("gRPC stopped.")
}
