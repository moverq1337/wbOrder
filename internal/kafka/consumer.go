package kafka

import (
	"context"
	"encoding/json"
	"github.com/moverq1337/wbOrder/internal/app/redis"
	"github.com/moverq1337/wbOrder/internal/models"
	"github.com/segmentio/kafka-go"
	"gorm.io/gorm"
	"log"
	"time"
)

type ConsumerConfig struct {
	Brokers []string
	Topic   string
	GroupID string
	DB      *gorm.DB
}

func NewConsumer(cfg ConsumerConfig) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  cfg.Brokers,
		Topic:    cfg.Topic,
		GroupID:  cfg.GroupID,
		MaxBytes: 10e6,
	})
}

func ConsumeMessages(ctx context.Context, reader *kafka.Reader, db *gorm.DB) {
	if redis.Client == nil {
		log.Println("Redis недоступен, работаем без кэша")
	}

	for {
		select {
		case <-ctx.Done():
			log.Println("Получен сигнал остановки Kafka consumer")
			return
		default:
			msg, err := reader.ReadMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					return
				}
				log.Printf("Ошибка чтения: %v", err)
				continue
			}

			if !json.Valid(msg.Value) {
				log.Printf("Невалидный JSON: %s", string(msg.Value))
				continue
			}

			var order models.Order
			if err = json.Unmarshal(msg.Value, &order); err != nil {
				log.Printf("Ошибка парсинга: %v", err)
				continue
			}

			tx := db.Begin()
			if err = tx.Create(&order).Error; err != nil {
				tx.Rollback()
				log.Printf("Ошибка сохранения %s: %v", order.OrderUID, err)
				continue
			}
			tx.Commit()

			orderJSON, err := json.Marshal(order)
			if err != nil {
				log.Printf("Ошибка сериализации %s: %v", order.OrderUID, err)
				continue
			}

			if redis.Client != nil {
				err = redis.Client.Set(context.Background(), order.OrderUID, orderJSON, time.Minute*10).Err()
				if err != nil {
					log.Printf("Ошибка Redis %s: %v", order.OrderUID, err)
				}
			}
		}
	}
}
