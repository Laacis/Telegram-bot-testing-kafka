package order_generation_service

import (
	"bufio"
	"fmt"
	models "order_generation_service/models"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func LoadSQLFile[T any](filename string, tableName string, storage *InMemoryStorage[T], parseFunc func([]string) (T, error), pattern string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	re := regexp.MustCompile(pattern)
	var capture bool
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, "INSERT INTO") && strings.Contains(line, tableName) {
			capture = true
			continue
		}

		if capture {
			if strings.Contains(line, ";") {
				parseStatement(line, storage, parseFunc, re)
				capture = false
			} else {
				parseStatement(line, storage, parseFunc, re)
			}
		}
	}
	return scanner.Err()
}

func parseStatement[T any](line string, storage *InMemoryStorage[T], parseFunc func([]string) (T, error), re *regexp.Regexp) {

	matches := re.FindStringSubmatch(line)
	if matches != nil {
		parts := matches[1:]

		item, err := parseFunc(parts)
		if err == nil {
			storage.Add(item)
		} else {
			fmt.Printf("Error parsing parts: %v\n", err)
		}
	}
}

func ParseDestination(parts []string) (models.Destination, error) {
	customerID, err := strconv.Atoi(parts[4])
	if err != nil {
		return models.Destination{}, err
	}
	return models.Destination{
		RestaurantCode: parts[0],
		RestaurantName: parts[1],
		Address:        parts[2],
		AreaCode:       parts[3],
		CustomerId:     customerID,
	}, nil
}

func ParseProduct(parts []string) (models.Product, error) {
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
