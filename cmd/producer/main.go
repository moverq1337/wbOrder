package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	w := &kafka.Writer{
		Addr:     kafka.TCP("localhost:9092"),
		Topic:    "orders",
		Balancer: &kafka.LeastBytes{},
	}

	baseMessageMap := map[string]interface{}{
		"track_number": "TRACK001",
		"entry":        "WBIL",
		"delivery": map[string]interface{}{
			"name":    "Test User",
			"phone":   "+1234567890",
			"zip":     "2639809",
			"city":    "Kiryat Mozkin",
			"address": "Ploshad Mira 15",
			"region":  "Kraiot",
			"email":   "test@gmail.com",
		},
		"payment": map[string]interface{}{
			"transaction":   "pay123",
			"request_id":    "",
			"currency":      "USD",
			"provider":      "wbpay",
			"amount":        1817,
			"payment_dt":    1637907727,
			"bank":          "alpha",
			"delivery_cost": 1500,
			"goods_total":   317,
			"custom_fee":    0,
		},
		"items": []map[string]interface{}{
			{
				"chrt_id":      9934930,
				"track_number": "WBILMTESTTRACK",
				"price":        453,
				"rid":          "ab4219087a764ae0btest",
				"name":         "Mascaras",
				"sale":         30,
				"size":         "0",
				"total_price":  317,
				"nm_id":        2389212,
				"brand":        "Vivienne Sabo",
				"status":       202,
			},
		},
		"locale":             "en",
		"internal_signature": "",
		"customer_id":        "test",
		"delivery_service":   "meest",
		"shardkey":           "9",
		"sm_id":              99,
		"oof_shard":          "1",
	}

	sChan := make(chan os.Signal, 1)
	signal.Notify(sChan, syscall.SIGINT, syscall.SIGTERM)

	counter := 0

	for {
		orderUID := fmt.Sprintf("test%d_%d", time.Now().UnixNano(), counter)
		counter++

		messageMap := make(map[string]interface{})
		for k, v := range baseMessageMap {
			messageMap[k] = v
		}

		messageMap["date_created"] = time.Now().UTC().Format(time.RFC3339)
		messageMap["order_uid"] = orderUID
		messageMap["payment"].(map[string]interface{})["transaction"] = orderUID

		messageBytes, err := json.Marshal(messageMap)
		if err != nil {
			fmt.Printf("Ошибка JSON: %v\n", err)
			return
		}

		err = w.WriteMessages(context.Background(),
			kafka.Message{
				Key:   []byte("order_uid"),
				Value: messageBytes,
			})
		if err != nil {
			fmt.Printf("Ошибка отправк: %v\n", err)
		} else {
			fmt.Printf("Сообщение отправил! order_uid: %s\n", orderUID)
		}

		time.Sleep(2 * time.Second)

		select {
		case <-sChan:
			fmt.Println("Остановка...")
			w.Close()
			return
		default:
		}
	}
}
