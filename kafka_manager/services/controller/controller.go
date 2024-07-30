package kafka_manager

import (
	"bytes"
	"fmt"
	"io"
	producer "kafka_manager/services/producer"
	"net/http"
	"sync"
)

type KafkaController struct {
	producer *producer.KafkaProducer
	running  bool
	mu       sync.Mutex
}

func NewKafkaController() *KafkaController {
	return &KafkaController{
		producer: nil,
		running:  false,
	}
}

func (kc *KafkaController) ProducerSendMessagesHandler(writer http.ResponseWriter, request *http.Request) {

	//TODO check conditions
	if !kc.running {
		http.Error(writer, "Producer is not running", http.StatusInternalServerError)
		return
	}

	body, err := io.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, "invalid request payload", http.StatusBadRequest)
		return
	}
	messagesToSend := bytes.Split(body, []byte("\n"))

	//TODO make topic from  a source like .env or something else
	//TODO make new struct to track the status of the message and implement resending messages if status is not "sent : true"

	for _, message := range messagesToSend {
		err := kc.producer.SendBytes(message, "test topic")
		if err != nil {
			//TODO implement resend
			http.Error(writer, "Error sending message to Kafka: ", http.StatusInternalServerError)
			return
		}
	}
	writer.Write([]byte("All messages sent"))
}

func (kc *KafkaController) ProducerUpHandler(writer http.ResponseWriter, request *http.Request) {
	kc.mu.Lock()
	defer kc.mu.Unlock()
	fmt.Println(" Exec producer UP command ...")
	if kc.running {
		http.Error(writer, "Error: Producer already running", http.StatusBadRequest)
		return
	}

	brokers := []string{"kafka:9092"}
	prod, err := producer.NewKafkaProducer(brokers)
	if err != nil {
		http.Error(writer, "Error starting producer: "+err.Error(), http.StatusInternalServerError)
		return
	}
	kc.producer = prod
	kc.running = true
	writer.Write([]byte("Producer is up"))
}

func (kc *KafkaController) ProducerDownHandler(writer http.ResponseWriter, request *http.Request) {
	kc.mu.Lock()
	defer kc.mu.Unlock()

	fmt.Println(" Exec producer Down command ...")
	if !kc.running {
		fmt.Println("running is false")
		http.Error(writer, "Producer not running", http.StatusBadRequest)
		return
	}
	fmt.Println(" Exec Close() on producer")
	if err := kc.producer.Close(); err != nil {
		http.Error(writer, "Error stopping producer", http.StatusInternalServerError)
		return
	}
	fmt.Println(" go this far!!!")
	kc.producer = nil
	kc.running = false
	writer.Write([]byte("Producer stopped"))
}

func (kc *KafkaController) ProducerStatusHandler(writer http.ResponseWriter, request *http.Request) {
	fmt.Println(" Exec producer Status command ...")
	fmt.Printf("status: running is %t and producer exists: %t", kc.running, kc.producer != nil)
	var response string
	if kc.running {
		response = fmt.Sprintf("Producer is up and running")
	} else {
		response = fmt.Sprintf("Producer is down")
	}

	writer.Write([]byte(response))
}
