package order_generation_service

import (
	"database/sql"
	"fmt"
	"log"
	models "order_generation_service/models"
	"os"
	"strings"
)

type Customer = models.Customer
type Product = models.Product
type Destination = models.Destination

var customersDbPrefix = os.Getenv("CUSTOMERS_DB_NAME")
var warehouseDbPrefix = os.Getenv("WAREHOUSE_DB_NAME")

func FetchProducts() (*[]Product, error) {
	db, err := openDbConnection(warehouseDbPrefix)
	defer func() { _ = db.Close() }()

	products, err := products(db)
	if err != nil {
		log.Fatalf("Error retrieving products: %v", err)
	}
	return &products, nil
}

func FetchCustomers() (*[]Customer, error) {
	db, err := openDbConnection(customersDbPrefix)
	defer func() { _ = db.Close() }()

	customers, err := customers(db)
	if err != nil {
		log.Fatalf("Error retrieving customers: %v", err)
	}

	return &customers, nil
}

func FetchDestinations() (*[]Destination, error) {
	db, err := openDbConnection(customersDbPrefix)

	defer func() { _ = db.Close() }()

	destinations, err := destinations(db)
	if err != nil {
		log.Fatalf("Error retrieving destinations: %v", err)
	}
	return &destinations, nil
}

func openDbConnection(dbPrefix string) (*sql.DB, error) {
	dbUserString := fmt.Sprintf("%s_USER", dbPrefix)
	dbUser := os.Getenv(dbUserString)
	dbPasswordString := fmt.Sprintf("%s_PASSWORD", dbPrefix)
	dbPassword := os.Getenv(dbPasswordString)
	dbNameString := fmt.Sprintf("%s_NAME", dbPrefix)
	dbName := os.Getenv(dbNameString)
	dbHost := strings.ToLower(dbPrefix)
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=disable", dbUser, dbPassword, dbName, dbHost)
	fmt.Println(connStr)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	return db, err
}

func products(db *sql.DB) ([]Product, error) {
	rows, err := db.Query("SELECT product_key, name, manufacturer, thermal_category, buy_price, units_per_pallet, unit_weight_kg FROM products")
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var products []Product
	for rows.Next() {
		var product Product
		err := rows.Scan(&product.ProductKey, &product.Name, &product.Manufacturer, &product.ThermalCategory, &product.BuyPrice, &product.UnitsPerPallet, &product.UnitWeightKg)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}

func customers(db *sql.DB) ([]Customer, error) {
	rows, err := db.Query("SELECT id, restaurant_id, name, contact_number, tax_number, address FROM customers")
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

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

func destinations(db *sql.DB) ([]Destination, error) {
	rows, err := db.Query("SELECT id, restaurant_code, restaurant_name, address, area_code, customer_id FROM  destinations")
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

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
