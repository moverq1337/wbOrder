package db

import (
	"fmt"
	models2 "github.com/moverq1337/wbOrder/internal/models"
)

func Migration() {
	fmt.Println("Starting migration...")
	err := DataBase.AutoMigrate(&models2.Order{}, &models2.Item{}, &models2.Payment{}, &models2.Delivery{})
	if err != nil {
		fmt.Println("AutoMigrate failed:", err)
		return
	}
	fmt.Println("AutoMigrate successful")
}
