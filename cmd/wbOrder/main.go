package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/moverq1337/wbOrder/internal/app/db"
	"github.com/moverq1337/wbOrder/internal/app/redis"
	"github.com/moverq1337/wbOrder/internal/kafka"
)

var migrate = flag.Bool("migrate", false, "Run database migration")

func main() {
	flag.Parse()
	db.Connection()

	if err := redis.Connect(); err != nil {
		fmt.Printf("Редис не завелся, продолжим без кэша", err)
	}

	if *migrate {
		db.Migration()
	} else {
		fmt.Println("Migration skipped, because u run it with out flag (-migration)")
	}

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

	ctx := context.Background()
	if err := service.Start(ctx); err != nil {
		fmt.Printf("Ошибка запуска: %v\n", err)
	}

}
