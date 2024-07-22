package order_generation_service

import (
	models "order_generation_service/models"
)

// make a local storage to store []Order
type Order = models.Order

type OrderStorageInterface interface {
	AddOrder(order Order)
	NextOrder() (Order, bool)
	Length() int
}
type OrderStorage struct {
	orders []Order
	Length int
}

func NewStorage() *OrderStorage {
	return &OrderStorage{
		orders: make([]Order, 0),
		Length: 0,
	}
}

func (s *OrderStorage) AddOrder(order Order) {
	s.orders = append(s.orders, order)
	s.Length++
}

func (s *OrderStorage) NextOrder() (Order, bool) {
	if len(s.orders) == 0 {
		return Order{}, false
	}
	nextOrder := s.orders[0]
	s.orders = s.orders[1:]
	s.Length--
	return nextOrder, true
}
