package kafka_manager

import "github.com/IBM/sarama"

type KafkaProducer struct {
	producer sarama.SyncProducer
}

func NewKafkaProducer(brokers []string) (*KafkaProducer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}
	return &KafkaProducer{producer: producer}, nil
}

func (prod *KafkaProducer) SendMessage(topic, message string) error {
	messageToSend := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}
	_, _, err := prod.producer.SendMessage(messageToSend)
	return err
}

func (prod *KafkaProducer) Close() error {
	return prod.producer.Close()
}

func (prod *KafkaProducer) SendBytes(message []byte, topic string) interface{} {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(message),
	}

	_, _, err := prod.producer.SendMessage(msg)
	return err
}
