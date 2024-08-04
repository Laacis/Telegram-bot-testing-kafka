package main

import (
	"fmt"
	"github.com/gorilla/mux"
	controller "kafka_manager/services/controller"
	"log"
	"net/http"
)

type kc = controller.KafkaController

func main() {
	kafkaController := controller.NewKafkaController()
	fmt.Sprintln("Running manager")
	// when GET on /producer/up - run producer service
	// when POST on /producer/feed - receive message for kafka producer to post to kafka
	// when GET on /producer/down - stop producer service
	router := mux.NewRouter()
	router.HandleFunc("/producer/up", kafkaController.ProducerUpHandler).Methods("GET")
	router.HandleFunc("/producer/feed", kafkaController.ProducerSendMessagesHandler).Methods("POST")
	router.HandleFunc("/producer/down", kafkaController.ProducerDownHandler).Methods("GET")
	router.HandleFunc("/producer/status", kafkaController.ProducerStatusHandler).Methods("GET")
	http.Handle("/", router)
	log.Fatal(http.ListenAndServe(":8082", nil))
}
