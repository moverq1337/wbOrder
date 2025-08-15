package main

import (
	"context"
	"flag"
	"github.com/gin-contrib/cors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/moverq1337/wbOrder/internal/app/db"
	"github.com/moverq1337/wbOrder/internal/app/redis"
	http2 "github.com/moverq1337/wbOrder/internal/http"
	"github.com/moverq1337/wbOrder/internal/kafka"
)

var migrate = flag.Bool("migrate", false, "Run database migration")

func main() {
	flag.Parse()

	db.Connection()
	if db.DataBase == nil {
		log.Fatal("Ошибка: подключение к БД")
	}

	if err := redis.Connect(); err != nil {
		log.Printf("Редис не завелся, продолжим без кэша", err)
	}

	if *migrate {
		db.Migration()
	} else {
		log.Println("Migration skipped, because u run it with out flag (-migration)")
	}

	// Контекст для graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	// Запуск HTTP сервера в горутине
	wg.Add(1)
	go func() {
		defer wg.Done()
		startGinServer(ctx)
	}()

	// Запуск Kafka consumer в горутине
	wg.Add(1)
	go func() {
		defer wg.Done()
		startKafkaConsumer(ctx)
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Получен сигнал завершения")

	cancel()

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		log.Println("Все компоненты завершены")
	case <-time.After(30 * time.Second):
		log.Println("Таймаут graceful shutdown")
	}
}

func startGinServer(ctx context.Context) {
	r := gin.Default()
	r.Use(cors.Default())
	r.GET("/order/:orderUID", http2.OrderHandler)

	srv := &http.Server{
		Addr:    ":8081",
		Handler: r,
	}

	go func() {
		<-ctx.Done()
		log.Println("Остановка HTTP сервера...")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Printf("Ошибка остановки HTTP сервера: %v", err)
		}
	}()

	log.Println("HTTP сервер запущен на :8081")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Printf("Ошибка HTTP сервера: %v", err)
	}
}

func startKafkaConsumer(ctx context.Context) {
	kafkaCfg := kafka.ConsumerConfig{
		Brokers: []string{"localhost:9091", "localhost:9092", "localhost:9093"},
		Topic:   "orders",
		GroupID: "order-group",
		DB:      db.DataBase,
	}

	serviceCfg := kafka.ServiceConfig{
		ConsumerConfig: kafkaCfg,
		DB:             db.DataBase,
	}

	service := kafka.NewService(serviceCfg)

	log.Println("Запуск Kafka consumer...")
	if err := service.Start(ctx); err != nil {
		if ctx.Err() == nil {
			log.Printf("Ошибка Kafka consumer: %v", err)
		}
	}
	log.Println("Kafka consumer остановлен")
}
