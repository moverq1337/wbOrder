package models

import "time"

type Order struct {
	OrderUID          string    `gorm:"primaryKey" json:"order_uid"`
	TrackNumber       string    `json:"track_number"`
	Entry             string    `json:"entry"`
	Delivery          Delivery  `gorm:"foreignKey:OrderID" json:"delivery"`
	Payment           Payment   `gorm:"foreignKey:OrderID" json:"payment"`
	Items             []Item    `gorm:"foreignKey:OrderID" json:"items"`
	Locale            string    `json:"locale"`
	InternalSignature string    `json:"internal_signature"`
	CustomerID        string    `json:"customer_id"`
	DeliveryService   string    `json:"delivery_service"`
	Shardkey          string    `json:"shardkey"`
	SmID              int       `json:"sm_id"`
	DateCreated       time.Time `json:"date_created"`
	OofShard          string    `json:"oof_shard"`
}
