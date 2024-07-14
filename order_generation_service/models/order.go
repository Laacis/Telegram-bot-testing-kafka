package order_generation_service

import (
	"github.com/google/uuid"
	"time"
)

type Order struct {
	OrderId      uuid.UUID `json:"order_id"`
	CreationTime time.Time `json:"creation_time"`
	CustomerId   uuid.UUID `json:"customer_id"`
	Products     []string  `json:"products"`
	//Products        []Product `json:"products"`
	Delivery bool `json:"delivery"`
	//DeliveryAddress Address   `json:"delivery_address"`
}
