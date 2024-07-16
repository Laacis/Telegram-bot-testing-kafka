package order_generation_service

type Customer struct {
	ID            int    `json:"customer_id"`
	RestaurantID  string `json:"restaurant_id"`
	Name          string `json:"name"`
	ContactNumber string `json:"contact_number"`
	TaxNumber     string `json:"tax_number"`
	Address       string `json:"address"`
}
