package telegram_bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
	"strconv"
)

const (
	botTimeout                   = 60
	DebugBot                     = true
	producerManagerServicePrefix = "http://kafka-manager:8082"
	orderServicePrefix           = "http://order-service:8081"
)

var (
	kafkaProducerUpEndpoint      = producerManagerServicePrefix + "/producer/up"
	kafkaProducerDownEndpoint    = producerManagerServicePrefix + "/producer/down"
	kafkaProducerStatusEndpoint  = producerManagerServicePrefix + "/producer/status"
	orderServiceGenerateEndpoint = orderServicePrefix + "/generate-orders/"
	orderServiceSendAllEndpoint  = orderServicePrefix + "/orders/send/all"
	orderServiceSendEndpoint     = orderServicePrefix + "/orders/send/"
)

func BotToken() string {
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		log.Fatalf("TELEGRAM_BOT_TOKEN must be set")
	}
	return botToken
}

func BotUpdateConfig() tgbotapi.UpdateConfig {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = botTimeout
	return u
}

func GetEndpoint(command string, i ...int) (string, error) {
	var endpoint string
	switch command {
	case "producerUp":
		endpoint = kafkaProducerUpEndpoint
	case "producerDown":
		endpoint = kafkaProducerDownEndpoint
	case "producerStatus":
		endpoint = kafkaProducerStatusEndpoint
	case "sendAll":
		endpoint = orderServiceSendAllEndpoint
	case "send":
		if len(i) == 0 {
			return "", fmt.Errorf("missing parameter for /send command")
		}
		endpoint = orderServiceSendEndpoint + strconv.Itoa(i[0])
	case "generate":
		if len(i) == 0 {
			return "", fmt.Errorf("missing parameter for /generate command")
		}
		endpoint = orderServiceGenerateEndpoint + strconv.Itoa(i[0])
	default:
		return "", fmt.Errorf("unknown command: %s", command)
	}
	return endpoint, nil
}
