package order_generation_service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	models "order_generation_service/models"
	database "order_generation_service/services/database"
	generator "order_generation_service/services/generator"
	inmemory "order_generation_service/services/storage/inmemory"
	"strconv"
)

const (
	kafkaFeedUrl      = "http://kafka_manager:8082/producer/feed"
	messageChunkSize  = 100
	defaultOrderLimit = -1
)

var useInMemory bool

type Destination = models.Destination
type Product = models.Product
type Order = models.Order
type Response = models.Response

type Handler struct {
	Storage      *inmemory.Queue[Order]
	Destinations *inmemory.InMemoryStorage[Destination]
	Products     *inmemory.InMemoryStorage[Product]
}

func SetInMemoryUse(b bool) {
	useInMemory = b
}
func (h *Handler) GenerateOrdersHandler(writer http.ResponseWriter, request *http.Request) {
	numberOfOrders, ok := firstArgument(writer, request)
	if !ok {
		http.Error(writer, "Error generating orders", http.StatusInternalServerError)
		return
	}

	destinations, products, err := h.fetchData(writer)
	if okDestinations := verifyData(destinations); !okDestinations {
		http.Error(writer, "Error verifying Destinations", http.StatusInternalServerError)
		return
	}

	if okProducts := verifyData(products); !okProducts {
		http.Error(writer, "Error verifying Products", http.StatusInternalServerError)
		return
	}

	orders, err := generator.GenerateOrders(destinations, products, numberOfOrders)
	if err != nil {
		http.Error(writer, "Error generating orders", http.StatusInternalServerError)
		return
	}
	var counter int
	for _, order := range *orders {
		if !h.Storage.IsFull() {
			h.Storage.Enqueue(order)
			counter++
		}
	}
	response := Response{
		ResponseCode: http.StatusOK,
		Message:      fmt.Sprintf("Report: %d orders generated", counter),
	}
	marshalResponse, _ := json.Marshal(&response)
	writer.Write(marshalResponse)
}

func (h *Handler) SendOrdersHandler(writer http.ResponseWriter, request *http.Request) {
	if h.Storage.IsEmpty() {
		http.Error(writer, "Storage empty, no orders to send", http.StatusInternalServerError)
		return
	}
	var messages [][]byte
	var counter int
	ordersToSend, ok := firstArgument(writer, request)
	if !ok {
		messages, counter = h.prepareOrdersToSend(writer)
	} else {
		messages, counter = h.prepareOrdersToSend(writer, ordersToSend)
	}

	for _, message := range messages {
		req, err := http.NewRequest(http.MethodPost, kafkaFeedUrl, bytes.NewReader(message))
		if err != nil {
			http.Error(writer, "Error crafting http request"+err.Error(), http.StatusInternalServerError)
			return
		}

		client := &http.Client{}
		response, err := client.Do(req)
		if err != nil {
			http.Error(writer, "Error sending request"+err.Error(), http.StatusInternalServerError)
			return
		}

		if response.StatusCode != http.StatusOK {
			http.Error(writer, "Error in response from kafka producer service"+err.Error(), http.StatusInternalServerError)
			return
		}
	}
	response := Response{
		ResponseCode: http.StatusOK,
		Message:      fmt.Sprintf("Report: successfully sent %d orders", counter),
	}
	f, _ := json.Marshal(&response)
	writer.Write(f)
}

func firstArgument(writer http.ResponseWriter, request *http.Request) (int, bool) {
	vars := mux.Vars(request)
	i, err := strconv.Atoi(vars["i"])
	if err != nil || i <= 0 {
		http.Error(writer, "Invalid parameter value", http.StatusBadRequest)
		return 0, false
	}
	return i, true
}

func (h *Handler) prepareOrdersToSend(writer http.ResponseWriter, ordersToSend ...int) ([][]byte, int) {
	var message []byte
	var counter int
	var messages [][]byte
	var messagesTotal int
	var messageLimit int
	if len(ordersToSend) > 0 && ordersToSend[0] >= 0 {
		messageLimit = ordersToSend[0]
	} else {
		messageLimit = defaultOrderLimit
	}

	for {
		nextOrder, more := h.Storage.Dequeue()
		if !more || messageLimit == 0 {
			messages = append(messages, message)
			break
		}
		if len(message) != 0 {
			message = append(message, '\n')
		}
		crafted, err := json.Marshal(nextOrder)
		if err != nil {
			http.Error(writer, "Error marshalling order", http.StatusInternalServerError)
			return nil, 0
		}
		message = append(message, crafted...)
		counter++
		messagesTotal++
		messageLimit--
		if counter == messageChunkSize {
			messages = append(messages, message)
			counter = 0
		}
	}
	return messages, messagesTotal
}
func verifyData[T any](objects *[]T) bool {

	return len(*objects) > 0
}

func (h *Handler) fetchData(writer http.ResponseWriter) (*[]Destination, *[]Product, error) {
	var destinations *[]Destination = nil
	var products *[]Product = nil
	var err error
	if useInMemory {
		destinations = h.Destinations.AllRecords()
		products = h.Products.AllRecords()
	} else {
		destinations, err = database.FetchDestinations()
		if err != nil {
			http.Error(writer, "Error fetching data from customers db", http.StatusInternalServerError)
		}

		products, err = database.FetchProducts()
		if err != nil {
			http.Error(writer, "Error fetching data from customers db", http.StatusInternalServerError)
		}
	}
	return destinations, products, err
}
