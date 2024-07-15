package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"order_generation_service/models"
	"os"
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
	s, _ := fetchData()
	//response := fmt.Sprintf("/%", s)
	writer.Write([]byte(s))
}

func fetchData() (string, error) {
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	dbHost := "db"
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=disable", dbUser, dbPassword, dbName, dbHost)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return "Successfully connected to the database!", nil
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
