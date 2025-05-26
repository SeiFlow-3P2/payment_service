package main

import (
    "context"
    "log"
    "os"
    "os/signal"
    "syscall"
    "time"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"

    "github.com/SeiFlow-3P2/payment_service/internal/app"
)

func main() {
    // 1. Чтение переменной окружения с подключением к БД
    dsn := os.Getenv("DATABASE_URL")
if dsn == "" {
    // временно для локальной отладки
    dsn = "postgres://postgres:12345@localhost:5432/payment?sslmode=disable"
}


    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalf("failed to connect to database: %v", err)
    }

    // 2. Создаём приложение
    application := &app.App{}

    // 3. Запуск серверов
    if err := application.Start(":50051", ":8080", db); err != nil {
        log.Fatalf("failed to start application: %v", err)
    }

    // 4. Ожидание завершения по сигналу (Ctrl+C или SIGTERM)
    stop := make(chan os.Signal, 1)
    signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

    <-stop
    log.Println("Termination signal received. Shutting down...")

    // 5. Graceful shutdown
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := application.Shutdown(ctx); err != nil {
        log.Fatalf("error during shutdown: %v", err)
    }

    log.Println("Application stopped.")
}
