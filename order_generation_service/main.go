package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	models "order_generation_service/models"
	database "order_generation_service/services/database"
	"strconv"
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
	router.HandleFunc("/fetch-products", fetchProductDataHandler)
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

func fetchProductDataHandler(w http.ResponseWriter, r *http.Request) {

	//Pseudo:
	//connect with customers db and retrieve customer data and destinations

	// connect to Warehouse db and retrieve product

	// Process the data and generate an order
	data, err := database.FetchProductData()
	if err != nil {
		log.Println("Error fetching data:", err)
		http.Error(w, "Error fetching data", http.StatusInternalServerError)
		return
	}
	w.Write([]byte(data))
}
