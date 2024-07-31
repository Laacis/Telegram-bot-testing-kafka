package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const (
	botTimeout                   = 60
	kafkaProducerUpEndpoint      = "http://kafka_manager:8082/producer/up"
	kafkaProducerDownEndpoint    = "http://kafka_manager:8082/producer/down"
	kafkaProducerStatusEndpoint  = "http://kafka_manager:8082/producer/status"
	orderServiceGenerateEndpoint = "http://order-service:8081/generate-orders/"
	orderServiceSendAllEndpoint  = "http://order-service:8081/orders/send/all"
)

var debugBot = true

func main() {
	loadEnv()
	botToken := token()
	bot := botSetup(botToken)
	u := botConfigUpdateTimeout()
	updates := bot.GetUpdatesChan(u)
	updateHandler(updates, bot)
}

func botConfigUpdateTimeout() tgbotapi.UpdateConfig {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = botTimeout
	return u
}

func botSetup(botToken string) *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = debugBot
	log.Printf("Authorized on account %s", bot.Self.UserName)
	return bot
}

func updateHandler(updates tgbotapi.UpdatesChannel, bot *tgbotapi.BotAPI) {
	for update := range updates {
		if update.Message == nil || !update.Message.IsCommand() {
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		command := update.Message.Command()
		var response string
		var err error

		switch command {
		case "help":
			response = handleHelp()
		case "producerUp", "producerDown", "producerStatus", "sendAll":
			response, err = commonCalls(command)
		case "generate":
			var numberOfOrders int
			command, numberOfOrders, err = commandIntoArguments(update.Message.Text)
			if err != nil {
				continue
			}
			response, err = commonCalls(command, numberOfOrders)
		case "status":
			response = "Server status: "
		default:
			response = "Unknown command"
		}

		if err != nil {
			log.Printf("Error handling command: %s: %v", command, err)
			response = "An error occurred while handling your request" + err.Error()
		}

		msg = tgbotapi.NewMessage(update.Message.Chat.ID, response)
		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}
}

func commandIntoArguments(message string) (string, int, error) {
	splitText := strings.Split(message, " ")
	command := splitText[0]
	ordersCount, err := strconv.Atoi(splitText[1])
	if err != nil {
		log.Println("Error converting str to int", err)
		return "", 0, err
	}
	return command, ordersCount, nil
}

func handleHelp() string {
	return "I understand /producerUp /producerDown /producerStatus '/generate n' and /status."
}

func token() string {
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		log.Fatalf("TELEGRAM_BOT_TOKEN must be set")
	}
	return botToken
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func commonCalls(command string, i ...int) (string, error) {
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
	case "generate":
		endpoint = orderServiceGenerateEndpoint + strconv.Itoa(i[0])
	}
	response, err := http.Get(endpoint)

	if err != nil {
		return "", err
	}
	defer func() { _ = response.Body.Close() }()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func ordersGenerate(i int) (string, error) {
	response, err := http.Get(orderServiceGenerateEndpoint + strconv.Itoa(i))
	if err != nil {
		return "", err
	}
	defer func() { _ = response.Body.Close() }()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
