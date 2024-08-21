package main

import (
	"flag"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	handler "order_generation_service/services/handler"
	storage "order_generation_service/services/storage"
	inmemory "order_generation_service/services/storage/in-memory"
)

var useInMemory bool
var orderStorage = storage.NewStorage()
var inMemoryDataBase *inmemory.InMemoryDataBase

func init() {
	flag.BoolVar(&useInMemory, "inmemory", false, "Use in-memory storage instead of database")
}

func main() {
	loadEnv()
	if useInMemory {
		inMemoryDataBase = inmemory.NewInMemoryDataBase()
		inMemoryDataBase.PopulateInMemoryStorage()
		handler.SetInMemoryUse(true)
	}

	handler := &handler.Handler{
		Storage:      orderStorage,
		Destinations: inMemoryDataBase.Destinations,
		Products:     inMemoryDataBase.Products,
	}

	router := mux.NewRouter()
	router.HandleFunc("/generate-orders/{i}", handler.GenerateOrdersHandler).Methods("GET")
	router.HandleFunc("/orders/send/all", handler.SendOrdersHandler).Methods("GET")
	router.HandleFunc("/orders/send/{i}", handler.SendOrdersHandler).Methods("GET")
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
