package order_generation_service

import (
	models "order_generation_service/models"
	"os"
	"testing"
)

// Test for loadSQLFile and parseStatement functions
func TestLoadSQLFile(t *testing.T) {
	// Create a temporary SQL file with test data
	sqlContent := `
		INSERT INTO destinations (restaurant_code, restaurant_name, address, area_code, customer_id) VALUES
		('MCD001', 'McD Victory park', '20A Victory park, Springfield, 5501', '55', 1),
		('KFC101', 'KFC SpringField W', '20B West st., Springfield, 5502', '55', 2),
		('JLB011', 'JolliBee SP. Trade Center', '5A Merchant gate, Springfield, 5505', '55', 3);
		
		INSERT INTO products (product_key, name, manufacturer, thermal_category, buy_price, units_per_pallet, unit_weight_kg) VALUES
		('893218913', 'Beef Pate S1', 'Crown Inc', 'FROZEN', 199.10, 40, 15),
		('712739324', 'Chicken Pate LX', 'Crown Inc', 'FROZEN', 189.50, 40, 15),
		('234567001', 'Classic Mayonnaise', 'Creamy Delights', 'CHILL', 45.99, 100, 1);
`

	tmpfile, err := os.CreateTemp("", "testdata*.sql")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tmpfile.Name()) // Clean up

	_, err = tmpfile.WriteString(sqlContent)
	if err != nil {
		t.Fatalf("Failed to write to temporary file: %v", err)
	}
	tmpfile.Close()

	//setting up regEx patterns
	destinationPattern := `\('([^']*)',\s*'([^']*)',\s*'([^']*)',\s*'([^']*)',\s*([0-9]+)\)`
	productPattern := `\('([^']*)',\s*'([^']*)',\s*'([^']*)',\s*'([^']*)',\s*([0-9]+(?:\.[0-9]+)?),\s*([0-9]+),\s*([0-9]+)\)`

	// Initialize in-memory storages
	storage := NewInMemoryStorage[models.Destination]()
	storageTwo := NewInMemoryStorage[models.Product]()

	// Load data from the temporary SQL file
	err = LoadSQLFile(tmpfile.Name(), "destinations", storage, ParseDestination, destinationPattern)
	if err != nil {
		t.Fatalf("Failed to load SQL file: %v", err)
	}

	// Load data from the temporary SQL file
	err = LoadSQLFile(tmpfile.Name(), "products", storageTwo, ParseProduct, productPattern)
	if err != nil {
		t.Fatalf("Failed to load SQL file: %v", err)
	}

	// Verify the contents of the in-memory storage
	expected := []models.Destination{
		{
			RestaurantCode: "MCD001",
			RestaurantName: "McD Victory park",
			Address:        "20A Victory park, Springfield, 5501",
			AreaCode:       "55",
			CustomerId:     1,
		},
		{
			RestaurantCode: "KFC101",
			RestaurantName: "KFC SpringField W",
			Address:        "20B West st., Springfield, 5502",
			AreaCode:       "55",
			CustomerId:     2,
		},
		{
			RestaurantCode: "JLB011",
			RestaurantName: "JolliBee SP. Trade Center",
			Address:        "5A Merchant gate, Springfield, 5505",
			AreaCode:       "55",
			CustomerId:     3,
		},
	}

	if storage.Length() != len(expected) {
		t.Fatalf("Expected %d records, got %d", len(expected), storage.Length())
	}

	for i, dest := range *storage.AllRecords() {
		if dest != expected[i] {
			t.Errorf("Expected %+v, got %+v", expected[i], dest)
		}
	}

	//Verify second storage
	if storageTwo.Length() != 3 {
		t.Fatalf("Expected %d records, got %d", 3, storageTwo.Length())
	}

	//TODO write down and complete compare expected items
}
