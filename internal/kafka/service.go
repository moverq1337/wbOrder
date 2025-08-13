package kafka

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/segmentio/kafka-go"
	"gorm.io/gorm"
)

type ServiceConfig struct {
	ConsumerConfig ConsumerConfig
	DB             *gorm.DB
}

type Service struct {
	consumer *kafka.Reader
	db       *gorm.DB
}

func NewService(cfg ServiceConfig) *Service {
	return &Service{
		consumer: NewConsumer(cfg.ConsumerConfig),
		db:       cfg.DB,
	}
}

func (s *Service) Start(ctx context.Context) error {
	go ConsumeMessages(ctx, s.consumer, s.db)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	fmt.Println("Остановка...")
	s.consumer.Close()
	return nil
}
