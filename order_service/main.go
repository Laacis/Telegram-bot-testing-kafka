package main

import (
	"bytes"
	"encoding/json"
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

var useInMemory bool
var keep = storage.NewStorage()
var destinationsInMemory = inmemory.NewInMemoryStorage[Destination]()
var productsInMemory = inmemory.NewInMemoryStorage[Product]()

const (
	destinationPattern       = `\('([^']*)',\s*'([^']*)',\s*'([^']*)',\s*'([^']*)',\s*([0-9]+)\)`
	productPattern           = `\('([^']*)',\s*'([^']*)',\s*'([^']*)',\s*'([^']*)',\s*([0-9]+(?:\.[0-9]+)?),\s*([0-9]+),\s*([0-9]+)\)`
	customersSqlInitFileName = "./sql/init.sql"
	warehouseSqlInitFileName = "./sql/init_warehouse.sql"
	kafkaFeedUrl             = "http://kafka_manager:8082/producer/feed"
)

func init() {
	flag.BoolVar(&useInMemory, "inmemory", false, "Use in-memory storage instead of database")
}

func main() {
	loadEnv()
	if useInMemory {
		populateInMemoryStorage()
	}

	router := mux.NewRouter()
	router.HandleFunc("/generate-orders/{i}", generateOrdersHandler).Methods("GET")
	router.HandleFunc("/orders/send/all", sendOrdersHandler).Methods("GET")
	http.Handle("/", router)
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	flag.Parse()
}

func populateInMemoryStorage() {
	if err := inmemory.LoadSQLFile(customersSqlInitFileName, "destinations", destinationsInMemory, inmemory.ParseDestination, destinationPattern); err != nil {
		log.Fatalf("Failed to load SQL file: %v", err)
	}
	if err := inmemory.LoadSQLFile(warehouseSqlInitFileName, "products", productsInMemory, inmemory.ParseProduct, productPattern); err != nil {
		log.Fatalf("Failed to load SQL file: %v", err)
	}
}

func sendOrdersHandler(writer http.ResponseWriter, _ *http.Request) {
	if keep.Length() == 0 {
		http.Error(writer, "Storage empty, no orders to send", http.StatusInternalServerError)
		return
	}

	var jsonBody []byte
	var counter int
	for {
		nextOrder, more := keep.NextOrder()
		if !more {
			break
		}
		if len(jsonBody) != 0 {
			jsonBody = append(jsonBody, '\n')
		}
		crafted, err := json.Marshal(nextOrder)
		if err != nil {
			http.Error(writer, "Error marshalling order", http.StatusInternalServerError)
			return
		}
		jsonBody = append(jsonBody, crafted...)
		counter++
	}

	req, err := http.NewRequest(http.MethodPost, kafkaFeedUrl, bytes.NewReader(jsonBody))
	if err != nil {
		http.Error(writer, "Error crafting http request"+err.Error(), http.StatusInternalServerError)
		return
	}

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		http.Error(writer, "Error sending request"+err.Error(), http.StatusInternalServerError)
		return
	}

	if response.StatusCode != http.StatusOK {
		http.Error(writer, "Error in response from kafka producer service"+err.Error(), http.StatusInternalServerError)
		return
	}

	f := fmt.Sprintf("Report: %d successfully sent", counter)
	writer.Write([]byte(f))
}

func generateOrdersHandler(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	i, err := strconv.Atoi(vars["i"])
	if err != nil || i <= 0 {
		http.Error(writer, "Invalid parameter value", http.StatusBadRequest)
		return
	}

	var destinations *[]Destination = nil
	var products *[]Product = nil

	if useInMemory {
		destinations = destinationsInMemory.AllRecords()
		products = productsInMemory.AllRecords()
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
	records := keep.Length()

	finalStr := fmt.Sprintf("Report: %d orders generated", records)
	writer.Write([]byte(finalStr))
}
