package main

import (
	"fmt"
	kafka_manager "kafka_manager/services/producer"
)

func main() {

	fmt.Sprintln("Running manager")
	// create a producer
	// when GET on /producer/up - run producer service
	// when POST on /producer/feed - receive message for kafka producer to post to kafka
	// when GET on /producer/down - stop producer service
	brokers := []string{"localhost:8091"}
	producer, _ := kafka_manager.NewKafkaProducer(brokers)
}
