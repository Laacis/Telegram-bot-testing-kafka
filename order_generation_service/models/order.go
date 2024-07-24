package order_generation_service

import (
	"github.com/google/uuid"
	"time"
)

type Order struct {
	OrderId         uuid.UUID `json:"order_id"`
	CreationTime    time.Time `json:"creation_time"`
	CustomerId      int       `json:"customer_id"`
	RestaurantCode  string    `json:"restaurant_code"`
	RestaurantName  string    `json:"restaurant_name"`
	Products        []Product `json:"products"`
	IsDelivery      bool      `json:"isDelivery"`
	DeliveryAddress string    `json:"delivery_address"`
}
