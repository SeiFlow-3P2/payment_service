package app

import (
    "context"
    "fmt"
    "log"
    "net"
    "net/http"
    "time"

    "google.golang.org/grpc"
    "gorm.io/gorm"

    internalapi "github.com/SeiFlow-3P2/payment_service/internal/api"
    "github.com/SeiFlow-3P2/payment_service/internal/repository"
    "github.com/SeiFlow-3P2/payment_service/internal/service"
)

type App struct {
    grpcServer     *grpc.Server
    httpServer     *http.Server
    paymentService *service.PaymentService
    shutdownChan   chan struct{}
}

func (a *App) Start(grpcAddr, httpAddr string, db *gorm.DB) error {
    a.shutdownChan = make(chan struct{})

    // Репозиторий и сервис
    paymentRepo := repository.NewPaymentRecordGorm(db)
    a.paymentService = service.NewPaymentService(paymentRepo)

    // gRPC сервер
    a.grpcServer = grpc.NewServer()
    lis, err := net.Listen("tcp", grpcAddr)
    if err != nil {
        return fmt.Errorf("failed to listen on gRPC: %w", err)
    }

    go func() {
        if err := a.grpcServer.Serve(lis); err != nil {
            log.Printf("gRPC server error: %v", err)
        }
    }()
    log.Printf("gRPC server started at %s", grpcAddr)

    // HTTP сервер и хендлер с каналом завершения
    handler := internalapi.NewWebhookHandler(a.paymentService, a.shutdownChan)
    mux := http.NewServeMux()
    mux.HandleFunc("/webhook", handler.HandleStripeWebhook)

    a.httpServer = &http.Server{
        Addr:    httpAddr,
        Handler: mux,
    }

    go func() {
        if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Printf("HTTP server error: %v", err)
        }
    }()
    log.Printf("HTTP server started at %s", httpAddr)

    // слушаем shutdown-сигнал от webhook
    go func() {
        <-a.shutdownChan
        log.Println("Shutdown signal received (from webhook)")

        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()
        if err := a.Shutdown(ctx); err != nil {
            log.Printf("error during shutdown: %v", err)
        }
    }()

    return nil
}

func (a *App) Shutdown(ctx context.Context) error {
    log.Println("Shutting down servers...")

    if a.httpServer != nil {
        if err := a.httpServer.Shutdown(ctx); err != nil {
            return fmt.Errorf("error shutting down HTTP server: %w", err)
        }
    }
    if a.grpcServer != nil {
        a.grpcServer.GracefulStop()
    }

    log.Println("Servers shut down gracefully.")
    return nil
}
