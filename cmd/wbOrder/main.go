package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/moverq1337/wbOrder/internal/app/db"
	"github.com/moverq1337/wbOrder/internal/kafka"
	"os"
	"os/signal"
	"syscall"
)

var migrate = flag.Bool("migrate", false, "Run database migration")

func main() {
	flag.Parse()
	db.Connection()

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
	consumer := kafka.NewConsumer(kafkaCfg)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go kafka.ConsumeMessages(ctx, consumer, db.DataBase)

	sChan := make(chan os.Signal, 1)
	signal.Notify(sChan, syscall.SIGINT, syscall.SIGTERM)
	<-sChan

	consumer.Close()
	fmt.Println("service stopped")

}
