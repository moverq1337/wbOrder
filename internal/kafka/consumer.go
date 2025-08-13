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
	if err := redis.Connect(); err != nil {
		log.Printf("Ошибка подключения к Redis: %v без кэша", err)
	}
	for {
		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("Consumer error(Ошибка чтения: %v):", err)
			continue
		}

		var order models.Order
		if err = json.Unmarshal(msg.Value, &order); err != nil {
			log.Println("Ошибка парсинга", err)
			continue
		}

		tx := db.Begin()
		if err = tx.Create(&order).Error; err != nil {
			tx.Rollback()
			log.Printf("Ошибка сохранения заказа %s: %v", order.OrderUID, err)
			continue
		}
		tx.Commit()

		orderJSON, err := json.Marshal(order)
		if err != nil {
			log.Printf("Ошибка сериализации заказа %s: %v", order.OrderUID, err)
			continue
		}

		if redis.Client != nil {
			err = redis.Client.Set(context.Background(), order.OrderUID, orderJSON, time.Hour).Err()
			if err != nil {
				log.Printf("Ошибка сохранения в redis %s: %v", order.OrderUID, err)
			} else {
				log.Printf("Заказ %s успешно прошел кэширование", order.OrderUID)
			}
		}

		log.Printf("Получено и сохр: key=%s, value=%s\n", string(msg.Key), string(msg.Value))
	}

}
