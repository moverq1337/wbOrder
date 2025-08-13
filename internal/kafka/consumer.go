package kafka

import (
	"context"
	"encoding/json"
	"github.com/moverq1337/wbOrder/internal/models"
	"github.com/segmentio/kafka-go"
	"gorm.io/gorm"
	"log"
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

		log.Printf("Получено: key=%s, value=%s\n", string(msg.Key), string(msg.Value))
	}

}
