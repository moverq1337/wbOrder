package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/moverq1337/wbOrder/internal/app/db"
	"github.com/moverq1337/wbOrder/internal/models"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
	"time"
)

var Client *redis.Client

func Connect() error {
	addr := os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT")

	Client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	if err := Client.Ping(context.Background()).Err(); err != nil {
		return fmt.Errorf("redis connection failed: %v", err)
	}
	fmt.Println("REdis done")
	return InitializeCache()
}

func InitializeCache() error {
	if db.DataBase == nil {
		return fmt.Errorf("БД не подключается")
	}
	var orders []models.Order
	if err := db.DataBase.Order("date_created DESC").Limit(20).Find(&orders).Error; err != nil {
		return fmt.Errorf("Ошибка загрузки заказов из БД: %v", err)
	}

	for _, order := range orders {
		orderJSON, err := json.Marshal(order)
		if err != nil {
			log.Printf("Ошибка сериализации заказа %s: %v", order.OrderUID, err)
			continue
		}
		if err = Client.Set(context.Background(), order.OrderUID, orderJSON, time.Minute*10).Err(); err != nil {
			log.Printf("Ошибка кэширования заказа %s: %v", order.OrderUID, err)
		}
	}
	fmt.Println("Кэш успешно ушел и готов из БД")
	return nil
}
