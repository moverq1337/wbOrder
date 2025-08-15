package http

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/moverq1337/wbOrder/internal/app/db"
	"github.com/moverq1337/wbOrder/internal/app/redis"
	"github.com/moverq1337/wbOrder/internal/models"
	"net/http"
	"strings"
	"time"
)

func OrderHandler(c *gin.Context) {
	orderUID := strings.TrimSpace(c.Param("orderUID"))
	if orderUID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "order_uid is required"})
		return
	}

	orderJSON, err := redis.Client.Get(context.Background(), orderUID).Result()
	if err == nil && orderJSON != "" {
		var order models.Order
		if err := json.Unmarshal([]byte(orderJSON), &order); err == nil {
			c.JSON(http.StatusOK, gin.H{"order": order})
			return
		}
	}

	var order models.Order
	if err := db.DataBase.
		Preload("Delivery").
		Preload("Payment").
		Preload("Items").
		Where("order_uid = ?", orderUID).
		First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "order not found"})
		return
	}

	orderBytes, err := json.Marshal(order)
	if err == nil {
		redis.Client.Set(context.Background(), orderUID, orderBytes, 10*time.Minute)
	}

	c.JSON(http.StatusOK, gin.H{"order": order})
}
