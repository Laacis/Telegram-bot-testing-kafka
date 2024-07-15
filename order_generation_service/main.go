package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"order_generation_service/models"
	"strconv"
	"time"
)

type Order = order_generation_service.Order

func main() {
	fmt.Println("Running order generator service v0.0.1")

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	router := mux.NewRouter()
	router.HandleFunc("/generate-order", generateOrderHandler)
	router.HandleFunc("/generate-orders/{i}", generateBulkOrders)
	http.Handle("/", router)
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func generateBulkOrders(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	i, err := strconv.Atoi(vars["i"])
	if err != nil || i <= 0 {
		http.Error(writer, "Invalid parameter", http.StatusBadRequest)
		return
	}
	response := fmt.Sprintf("hit the bulk generator @ /generate-orders/%d", i)
	writer.Write([]byte(response))
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
