package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"net/http"
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

func main() {
	fmt.Println("Running order generator service v0.0.1")

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	http.HandleFunc("/generate-order", generateOrderHandler)
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func generateOrderHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("hit the spot!")
	data, err := fetchDataFromDB()
	if err != nil {
		log.Println("Error fetching data:", err)
		http.Error(w, "Error fetching data", http.StatusInternalServerError)
		return
	}

	// Process the data and generate an order
	order, _ := json.Marshal(data)

	w.Write(order)
}

func fetchDataFromDB() ([]Order, error) {
	// Simulated data fetching logic
	return []Order{
		{
			OrderId:      uuid.New(),
			CreationTime: time.Now(),
			CustomerId:   uuid.New(), //get from db
			Products:     []string{"product 1", "product 2"},
			Delivery:     false,
		},
		{
			OrderId:      uuid.New(),
			CreationTime: time.Now(),
			CustomerId:   uuid.New(), //get from db
			Products:     []string{"product 1", "product 3"},
			Delivery:     true,
		},
	}, nil
}
