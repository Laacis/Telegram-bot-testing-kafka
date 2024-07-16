package order_generation_service

import (
	"database/sql"
	"fmt"
	"log"
	models "order_generation_service/models"
	"os"
)

type Customer = models.Customer

type Destination = models.Destination

func FetchCustomerData() (string, error) {
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

	customers, err := getCustomers(db)
	if err != nil {
		log.Fatalf("Error retrieving customers: %v", err)
	}

	destinations, err := getDestinations(db)
	if err != nil {
		log.Fatalf("Error retrieving destinations: %v", err)
	}

	reportCustomers := fmt.Sprintf("Received %d customer details.\n", len(customers))
	reportDestinations := fmt.Sprintf("Received %d destination details.", len(destinations))

	return string(reportCustomers + reportDestinations), nil
}

func getCustomers(db *sql.DB) ([]Customer, error) {
	rows, err := db.Query("SELECT id, restaurant_id, name, contact_number, tax_number, address FROM customers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customers []Customer
	for rows.Next() {
		var customer Customer
		err := rows.Scan(&customer.ID, &customer.RestaurantID, &customer.Name, &customer.ContactNumber, &customer.TaxNumber, &customer.Address)
		if err != nil {
			return nil, err
		}
		customers = append(customers, customer)
	}

	return customers, nil
}

func getDestinations(db *sql.DB) ([]Destination, error) {
	rows, err := db.Query("SELECT id, restaurant_code, restaurant_name, address, area_code, customer_id FROM  destinations")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var destinations []Destination
	for rows.Next() {
		var destination Destination
		err := rows.Scan(&destination.Id, &destination.RestaurantCode, &destination.RestaurantName, &destination.Address, &destination.AreaCode, &destination.CustomerId)
		if err != nil {
			return nil, err
		}
		destinations = append(destinations, destination)
	}
	return destinations, nil
}
