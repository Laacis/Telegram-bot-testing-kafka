package telegram_bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
	"strconv"
)

var (
	DebugBot                     bool
	botTimeout                   int
	botToken                     string
	producerManagerServicePrefix string
	orderServicePrefix           string
	kafkaProducerUpEndpoint      string
	kafkaProducerDownEndpoint    string
	kafkaProducerStatusEndpoint  string
	orderServiceGenerateEndpoint string
	orderServiceSendAllEndpoint  string
	orderServiceSendEndpoint     string
)

func init() {
	//token
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		log.Fatalf("TELEGRAM_BOT_TOKEN must be set")
	}

	//timeout
	value, err := strconv.Atoi(os.Getenv("BOT_TIMEOUT"))
	if err != nil {
		log.Printf("Invalid BOT_TIMEOUT, defaulting to 60: %v", err)
		botTimeout = 60 // Default value
	} else {
		botTimeout = value
	}

	//debug
	DebugBot = os.Getenv("DEBUG_BOT") == "true"

	//server endpoint prefixes
	producerManagerServicePrefix = os.Getenv("PRODUCER_MANAGER_SERVICE_PREFIX")
	orderServicePrefix = os.Getenv("ORDER_SERVICE_PREFIX")

	//endpoints
	kafkaProducerUpEndpoint = producerManagerServicePrefix + "/producer/up"
	kafkaProducerDownEndpoint = producerManagerServicePrefix + "/producer/down"
	kafkaProducerStatusEndpoint = producerManagerServicePrefix + "/producer/status"
	orderServiceGenerateEndpoint = orderServicePrefix + "/generate-orders/"
	orderServiceSendAllEndpoint = orderServicePrefix + "/orders/send/all"
	orderServiceSendEndpoint = orderServicePrefix + "/orders/send/"
}

func BotToken() string {
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
