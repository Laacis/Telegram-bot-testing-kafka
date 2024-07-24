package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	models "order_generation_service/models"
	database "order_generation_service/services/database"
	generator "order_generation_service/services/generator"
	storage "order_generation_service/services/storage"
	inmemory "order_generation_service/services/storage/in-memory"
	"strconv"
)

type Destination = models.Destination
type Product = models.Product
type Customer = models.Customer

var keep = storage.NewStorage()
var useInMemory bool

// TODO make in-memory storages
var destinationsInMemory = inmemory.NewInMemoryStorage[Destination]()
var productsInMemory = inmemory.NewInMemoryStorage[Product]()

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
	if useInMemory {
		//populate inMemory storages

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

	var destinations *[]Destination = nil
	var products *[]Product = nil

	if useInMemory {
		destinations = inMemory.Destinations()
		products = inMemory.Products()
	} else {
		destinations, err = database.FetchDestinations()
		if err != nil {
			http.Error(writer, "Error fetching data from customers db", http.StatusInternalServerError)
		}
		products, err = database.FetchProductData()
		if err != nil {
			http.Error(writer, "Error fetching data from customers db", http.StatusInternalServerError)
		}
	}

	orders, err := generator.GenerateOrders(destinations, products, i)
	if err != nil {
		http.Error(writer, "Error generating orders", http.StatusInternalServerError)
	}

	for _, order := range *orders {
		keep.AddOrder(order)
	}

	report := fmt.Sprintf("Successfully generated %d orders.", keep.Length)
	writer.Write([]byte(report))
}
