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
	messageChunkSize         = 100
	defaultOrderLimit        = -1
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
	router.HandleFunc("/orders/send/{i}", sendOrdersHandler).Methods("GET")
	router.HandleFunc("/orders/stored", storedOrdersHandler).Methods("GET")
	http.Handle("/", router)
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func storedOrdersHandler(writer http.ResponseWriter, request *http.Request) {
	http.Error(writer, "not implemented method", http.StatusNotImplemented)
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

func sendOrdersHandler(writer http.ResponseWriter, request *http.Request) {
	if keep.Length() == 0 {
		http.Error(writer, "Storage empty, no orders to send", http.StatusInternalServerError)
		return
	}
	var messages [][]byte
	var counter int
	ordersToSend, ok := firstArgument(writer, request)
	if !ok {
		messages, counter = prepareOrdersToSend(writer)
	} else {
		messages, counter = prepareOrdersToSend(writer, ordersToSend)
	}

	for _, message := range messages {
		req, err := http.NewRequest(http.MethodPost, kafkaFeedUrl, bytes.NewReader(message))
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
	}

	f := fmt.Sprintf("Report: successfully sent %d orders", counter)
	writer.Write([]byte(f))
}

func prepareOrdersToSend(writer http.ResponseWriter, ordersToSend ...int) ([][]byte, int) {
	var message []byte
	var counter int
	var messages [][]byte
	var messagesTotal int
	var messageLimit int
	if len(ordersToSend) > 0 && ordersToSend[0] >= 0 {
		messageLimit = ordersToSend[0]
	} else {
		messageLimit = defaultOrderLimit
	}

	for {
		nextOrder, more := keep.NextOrder()
		if !more || messageLimit == 0 {
			messages = append(messages, message)
			break
		}
		if len(message) != 0 {
			message = append(message, '\n')
		}
		crafted, err := json.Marshal(nextOrder)
		if err != nil {
			http.Error(writer, "Error marshalling order", http.StatusInternalServerError)
			return nil, 0
		}
		message = append(message, crafted...)
		counter++
		messagesTotal++
		messageLimit--
		if counter == messageChunkSize {
			messages = append(messages, message)
			counter = 0
		}
	}
	return messages, messagesTotal
}

func generateOrdersHandler(writer http.ResponseWriter, request *http.Request) {
	numberOfOrders, ok := firstArgument(writer, request)
	if !ok {
		http.Error(writer, "Error generating orders", http.StatusInternalServerError)
		return
	}

	destinations, products, err := fetchData(writer)
	if okDestinations := verifyData(destinations); !okDestinations {
		http.Error(writer, "Error verifying Destinations", http.StatusInternalServerError)
		return
	}

	if okProducts := verifyData(products); !okProducts {
		http.Error(writer, "Error verifying Products", http.StatusInternalServerError)
		return
	}

	orders, err := generator.GenerateOrders(destinations, products, numberOfOrders)
	if err != nil {
		http.Error(writer, "Error generating orders", http.StatusInternalServerError)
		return
	}

	for _, order := range *orders {
		keep.AddOrder(order)
	}
	records := keep.Length()
	finalStr := fmt.Sprintf("Report: %d orders generated", records)
	writer.Write([]byte(finalStr))
}

func verifyData[T any](objects *[]T) bool {
	return len(*objects) > 0
}

func fetchData(writer http.ResponseWriter) (*[]Destination, *[]Product, error) {
	var destinations *[]Destination = nil
	var products *[]Product = nil
	var err error
	if useInMemory {
		destinations = destinationsInMemory.AllRecords()
		products = productsInMemory.AllRecords()
	} else {
		destinations, err = database.FetchDestinations()
		if err != nil {
			http.Error(writer, "Error fetching data from customers db", http.StatusInternalServerError)
		}

		products, err = database.FetchProducts()
		if err != nil {
			http.Error(writer, "Error fetching data from customers db", http.StatusInternalServerError)
		}
	}
	return destinations, products, err
}

func firstArgument(writer http.ResponseWriter, request *http.Request) (int, bool) {
	vars := mux.Vars(request)
	i, err := strconv.Atoi(vars["i"])
	if err != nil || i <= 0 {
		http.Error(writer, "Invalid parameter value", http.StatusBadRequest)
		return 0, false
	}
	return i, true
}
