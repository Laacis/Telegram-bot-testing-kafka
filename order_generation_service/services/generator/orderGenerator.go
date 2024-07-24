package order_generation_service

import (
	"github.com/google/uuid"
	"log"
	"math/rand"
	models "order_generation_service/models"
	"time"
)

type Customer = models.Customer
type Destination = models.Destination
type Product = models.Product
type Order = models.Order

const maxProducts = 10

func GenerateOrders(d *[]Destination, p *[]Product, numberOfOrders int) (*[]Order, error) {
	var orders []Order
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)

	for i := 0; i < numberOfOrders; i++ {
		x := random.Intn(len(*d))
		destination := (*d)[x]
		numberOfProducts := rand.Intn(maxProducts)
		products, err := products(p, numberOfProducts)
		if err != nil {
			log.Fatalf("Error generating products %v", err)
		}

		order := Order{
			OrderId:        uuid.New(),
			CreationTime:   time.Now(),
			CustomerId:     destination.CustomerId,
			RestaurantCode: destination.RestaurantCode,
			RestaurantName: destination.RestaurantName,
			Products:       products,
			IsDelivery:     i%5 != 0, //every 5th to be not delivery
		}
		orders = append(orders, order)
	}
	return &orders, nil
}

func billingAddress(c *[]Customer, id int) string {
	for _, customer := range *c {
		if customer.ID == id {
			return customer.Address
		}
	}
	return ""
}

func products(p *[]Product, i int) ([]Product, error) {
	var products []Product
	for j := 0; j < i; j++ {
		r := rand.Intn(len(*p))
		product := (*p)[r]
		products = append(products, product)
	}
	return products, nil
}
