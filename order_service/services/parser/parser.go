package order_generation_service

import (
	models "order_generation_service/models"
	"strconv"
)

type Destination = models.Destination
type Product = models.Product

type Parser[T any] interface {
	Parse(parts []string) (T, error)
}
type DestinationParser struct{}
type ProductParser struct{}

func ParseData[T any](parser Parser[T], parts []string) (T, error) {
	return parser.Parse(parts)
}

func (d DestinationParser) Parse(parts []string) (Destination, error) {
	customerID, err := strconv.Atoi(parts[4])
	if err != nil {
		return Destination{}, err
	}
	return models.Destination{
		RestaurantCode: parts[0],
		RestaurantName: parts[1],
		Address:        parts[2],
		AreaCode:       parts[3],
		CustomerId:     customerID,
	}, nil
}

func (p ProductParser) Parse(parts []string) (Product, error) {
	bPrice, _ := strconv.ParseFloat(parts[4], 64)
	unitCount, _ := strconv.Atoi(parts[5])
	unitWeight, _ := strconv.Atoi(parts[6])
	return models.Product{
		ProductKey:      parts[0],
		Name:            parts[1],
		Manufacturer:    parts[2],
		ThermalCategory: parts[3],
		BuyPrice:        bPrice,
		UnitsPerPallet:  unitCount,
		UnitWeightKg:    unitWeight,
	}, nil
}
