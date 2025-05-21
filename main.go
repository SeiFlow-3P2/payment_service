package main

import (
	"log"
	"net"
	"net/http"
	"os"

	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/SeiFlow-3P2/payment_service/internal/api"
	"github.com/SeiFlow-3P2/payment_service/internal/repository"
	"github.com/SeiFlow-3P2/payment_service/internal/service"
	pb "github.com/SeiFlow-3P2/payment_service/pkg/proto/v1"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	paymentRepo := repository.NewPaymentRecordGorm(db)
	paymentService := service.NewPaymentService(paymentRepo)
	paymentAPI := api.NewPaymentAPI(paymentService)

	go func() {
		grpcServer := grpc.NewServer()
		pb.RegisterPaymentServiceServer(grpcServer, paymentAPI)

		listener, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		log.Println("gRPC server listening on :50051")
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	httpHandler := api.NewWebhookHandler(paymentService)
	http.HandleFunc("/stripe/webhook", httpHandler.HandleStripeWebhook)

	log.Println("Serving Stripe webhook on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("failed to start HTTP server: %v", err)
	}
}
