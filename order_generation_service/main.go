package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	order_generation_service "order_generation_service/models"
	database "order_generation_service/services/database"
	generator "order_generation_service/services/generator"
	storage "order_generation_service/services/storage"
	"strconv"
)

type Destination = order_generation_service.Destination
type Product = order_generation_service.Product

var keep = storage.NewStorage()
var useInMemory bool

// TODO make in-memory storages
var customersInMemory []Destination
var warehouseInMemory []Product

func init() {
	flag.BoolVar(&useInMemory, "inmemory", false, "Use in-memory storage instead of database")
}

func main() {
	fmt.Println("Running order generator service v0.0.1")

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	flag.Parse()
	if useInMemory == true {
		customersInMemory = make([]Destination, 0)
		warehouseInMemory = make([]Product, 0)
	}

	router := mux.NewRouter()
	//router.HandleFunc("/fetch-products", fetchProductDataHandler)
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
	customers, destinations, err := database.FetchCustomerData()
	if err != nil {
		http.Error(writer, "Error fetching data from customers db", http.StatusInternalServerError)
	}
	products, err := database.FetchProductData()
	if err != nil {
		http.Error(writer, "Error fetching data from customers db", http.StatusInternalServerError)
	}

	orders, err := generator.GenerateOrders(customers, destinations, products, i)

	for _, order := range *orders {
		keep.AddOrder(order)
	}

	report := fmt.Sprintf("Successfully generated %d orders.", keep.Length)
	writer.Write([]byte(report))
}
