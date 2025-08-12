package main

import (
	"fmt"
	models2 "github.com/moverq1337/wbOrder/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=100.64.213.196 user=user password=admin dbname=wb port=5432 sslmode=disable TimeZone=Europe/Moscow"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("No base")
		return
	} else {
		fmt.Println("Connect successfull")
	}

	err = db.AutoMigrate(&models2.Order{}, &models2.Item{}, &models2.Payment{}, &models2.Delivery{})
	if err != nil {
		fmt.Println("AutoMigrate failed:", err)
		return
	}
	fmt.Println("AutoMigrate successful")
}
