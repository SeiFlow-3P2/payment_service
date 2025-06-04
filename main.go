package main

import (
    "context"
    "log"
    "net"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "fmt"
    "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
    "google.golang.org/grpc"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"

    "github.com/SeiFlow-3P2/payment_service/internal/api"
    "github.com/SeiFlow-3P2/payment_service/internal/config"
    "github.com/SeiFlow-3P2/payment_service/internal/models"
    "github.com/SeiFlow-3P2/payment_service/internal/repository"
    "github.com/SeiFlow-3P2/payment_service/internal/service"
    pb "github.com/SeiFlow-3P2/payment_service/pkg/proto/v1"
)

func main() {
    // Load configuration
    cfg, err := config.Load()
    if err != nil {
        log.Fatal("Failed to load configuration:", err)
    }

    // Initialize DB
    db, err := gorm.Open(postgres.Open(cfg.GetDSN()), &gorm.Config{})
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

    grpcAddr := fmt.Sprintf("%s:%d", cfg.Servers.GRPC.Host, cfg.Servers.GRPC.Port)
    grpcLis, err := net.Listen("tcp", grpcAddr)
    if err != nil {
        log.Fatalf("Failed to listen on gRPC: %v", err)
    }

    go func() {
        log.Printf("gRPC server started on %s", grpcAddr)
        if err := grpcServer.Serve(grpcLis); err != nil {
            log.Fatalf("gRPC server failed: %v", err)
        }
    }()

    // REST Gateway
    go func() {
        ctx := context.Background()
        ctx, cancel := context.WithCancel(ctx)
        defer cancel()

        mux := runtime.NewServeMux()
        opts := []grpc.DialOption{grpc.WithInsecure()}

        restAddr := fmt.Sprintf("%s:%d", cfg.Servers.REST.Host, cfg.Servers.REST.Port)
        grpcEndpoint := fmt.Sprintf("%s:%d", cfg.Servers.GRPC.Host, cfg.Servers.GRPC.Port)

        err := pb.RegisterPaymentServiceHandlerFromEndpoint(ctx, mux, grpcEndpoint, opts)
        if err != nil {
            log.Fatalf("Failed to register gRPC-Gateway: %v", err)
        }

        log.Printf("REST gateway started on %s", restAddr)
        if err := http.ListenAndServe(restAddr, mux); err != nil && err != http.ErrServerClosed {
            log.Fatalf("REST gateway failed: %v", err)
        }
    }()

    // Webhook endpoint
    webhookHandler := api.NewWebhookHandler(paymentService, make(chan struct{}))

    go func() {
        http.HandleFunc(cfg.Servers.Webhook.Path, webhookHandler.HandleStripeWebhook)
        webhookAddr := fmt.Sprintf("%s:%d", cfg.Servers.Webhook.Host, cfg.Servers.Webhook.Port)
        log.Printf("HTTP webhook server started on %s%s", webhookAddr, cfg.Servers.Webhook.Path)
        if err := http.ListenAndServe(webhookAddr, nil); err != nil && err != http.ErrServerClosed {
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
