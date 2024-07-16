package main

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	models "order_generation_service/models"
	database "order_generation_service/services/database"
	"strconv"
	"time"
)

type Order = models.Order
type Customer = models.Customer

func main() {
	fmt.Println("Running order generator service v0.0.1")

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	router := mux.NewRouter()
	router.HandleFunc("/generate-order", generateOrderHandler)
	router.HandleFunc("/generate-orders/{i}", generateOrdersHandler).Methods("GET")
	http.Handle("/", router)
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func generateOrdersHandler(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	i, err := strconv.Atoi(vars["i"])
	if err != nil || i <= 0 {
		http.Error(writer, "Invalid parameter", http.StatusBadRequest)
		return
	}
	s, _ := database.FetchCustomerData()
	//response := fmt.Sprintf("/%", s)
	writer.Write([]byte(s))
}

func generateOrderHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("hit the spot!")
	data, err := fetchDataFromDB()
	if err != nil {
		log.Println("Error fetching data:", err)
		http.Error(w, "Error fetching data", http.StatusInternalServerError)
		return
	}
	//Pseudo:
	//connect with customers db and retrieve customer data and destinations

	// connect to Warehouse db and retrieve product

	// Process the data and generate an order

	sumUp := fmt.Sprintf("Received %d customer details.", len(data))
	w.Write([]byte(sumUp))
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
