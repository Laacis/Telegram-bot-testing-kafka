package order_generation_service

type Destination struct {
	Id             int    `json:"id"`
	RestaurantCode string `json:"restaurantCode"`
	RestaurantName string `json:"restaurantName"`
	Address        string `json:"address"`
	AreaCode       string `json:"areaCode"`
	CustomerId     int    `json:"customerId"`
}
