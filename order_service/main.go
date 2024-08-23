package main

import (
	"flag"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	handler "order_generation_service/services/handler"
	inmemory "order_generation_service/services/storage/inmemory"
)

var useInMemory bool
var inMemoryDataBase *inmemory.InMemoryDataBase

func init() {
	flag.BoolVar(&useInMemory, "inmemory", false, "Use inmemory storage instead of database")
}

func main() {
	loadEnv()
	if useInMemory {
		inMemoryDataBase = inmemory.NewInMemoryDataBase()
		inMemoryDataBase.PopulateInMemoryStorage()
		handler.SetInMemoryUse(true)
	}

	handler := &handler.Handler{
		Storage:      inmemory.NewQueue[handler.Order](100),
		Destinations: inMemoryDataBase.Destinations,
		Products:     inMemoryDataBase.Products,
	}

	router := mux.NewRouter()
	router.HandleFunc("/generate-orders/{i}", handler.GenerateOrdersHandler).Methods("GET")
	// TODO SendAll is not working as expected and verify number of orders sent(counter works not as expected)
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
