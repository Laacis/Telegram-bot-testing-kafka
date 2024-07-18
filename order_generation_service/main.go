package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	database "order_generation_service/services/database"
	generator "order_generation_service/services/generator"
	"strconv"
)

func main() {
	fmt.Println("Running order generator service v0.0.1")

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
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
	report := fmt.Sprintf("Successfully generated %d orders.", len(*orders))
	//res, _ := json.Marshal(*orders)

	writer.Write([]byte(report))
}
