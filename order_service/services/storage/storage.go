package order_generation_service

import (
	models "order_generation_service/models"
)

type Order = models.Order

type OrderStorageInterface interface {
	AddOrder(order Order)
	NextOrder() (Order, bool)
	Length() int
}
type OrderStorage struct {
	orders []Order
}

func NewStorage() *OrderStorage {
	return &OrderStorage{
		orders: make([]Order, 0),
	}
}

func (s *OrderStorage) AddOrder(order Order) {
	s.orders = append(s.orders, order)
}

func (s *OrderStorage) NextOrder() (*Order, bool) {
	if len(s.orders) == 0 {
		return nil, false
	}
	nextOrder := s.orders[0]
	s.orders = s.orders[1:]
	return &nextOrder, true
}

func (s *OrderStorage) Length() int {
	return len(s.orders)
}
