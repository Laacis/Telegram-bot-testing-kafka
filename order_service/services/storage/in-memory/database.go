package order_generation_service

import (
	"bufio"
	"fmt"
	"log"
	models "order_generation_service/models"
	parsers "order_generation_service/services/parser"
	"os"
	"regexp"
	"strings"
)

const (
	destinationPattern       = `\('([^']*)',\s*'([^']*)',\s*'([^']*)',\s*'([^']*)',\s*([0-9]+)\)`
	productPattern           = `\('([^']*)',\s*'([^']*)',\s*'([^']*)',\s*'([^']*)',\s*([0-9]+(?:\.[0-9]+)?),\s*([0-9]+),\s*([0-9]+)\)`
	customersSqlInitFileName = "./sql/init.sql"
	warehouseSqlInitFileName = "./sql/init_warehouse.sql"
)

type Destination = models.Destination
type Product = models.Product

type InMemoryDataBase struct {
	Destinations *InMemoryStorage[Destination]
	Products     *InMemoryStorage[Product]
}

func NewInMemoryDataBase() *InMemoryDataBase {
	return &InMemoryDataBase{
		Destinations: NewInMemoryStorage[Destination](),
		Products:     NewInMemoryStorage[Product](),
	}
}

func (db *InMemoryDataBase) PopulateInMemoryStorage() {
	destinationParser := parsers.DestinationParser{}
	productParser := parsers.ProductParser{}
	if err := LoadSQLFile(customersSqlInitFileName, "destinations", db.Destinations, destinationParser, destinationPattern); err != nil {
		log.Fatalf("Failed to load SQL file: %v", err)
	}
	if err := LoadSQLFile(warehouseSqlInitFileName, "products", db.Products, productParser, productPattern); err != nil {
		log.Fatalf("Failed to load SQL file: %v", err)
	}
}

func LoadSQLFile[T any](filename string, tableName string, storage *InMemoryStorage[T], parser parsers.Parser[T], pattern string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()
	re := regexp.MustCompile(pattern)
	statements := scanForStatement(file, tableName)
	for _, statement := range statements {
		parseStatementToStorage(statement, storage, parser, re)
	}

	return nil
}

func scanForStatement(file *os.File, tableName string) []string {
	var capture bool
	scanner := bufio.NewScanner(file)
	statements := make([]string, 0)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, "INSERT INTO") && strings.Contains(line, tableName) {
			capture = true
			continue
		}

		if capture {
			if strings.Contains(line, ";") {
				capture = false
			}
			statements = append(statements, line)
		}
	}

	return statements
}

func parseStatementToStorage[T any](statement string, storage *InMemoryStorage[T], parser parsers.Parser[T], re *regexp.Regexp) {
	matches := re.FindStringSubmatch(statement)
	if matches != nil {
		parts := matches[1:]
		item, err := parsers.ParseData(parser, parts)
		if err == nil {
			storage.Add(item)
		} else {
			fmt.Printf("Error parsing parts: %v\n", err)
		}
	}
}
