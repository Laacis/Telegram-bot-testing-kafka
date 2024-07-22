package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func main() {
	fmt.Println("V1.0")
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		log.Fatalf("TELEGRAM_BOT_TOKEN must be set")
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message updates
			continue
		}

		if !update.Message.IsCommand() { // ignore any non-command Messages
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		// Extract the command from the Message.
		switch update.Message.Command() {
		case "help":
			msg.Text = "I understand /fetchProducts /gene5 and /status."
		case "fetchProducts":
			msg.Text = "executing generate orders..."
			// call order-generation-service
			orders, err := callOrderGenerationService()
			if err != nil {
				log.Println("Error generating order:", err)
				continue
			}
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, orders)
		case "generate":
			splitText := strings.Split(update.Message.Text, " ")
			ordersCount, err := strconv.Atoi(splitText[1])
			if err != nil {
				log.Println("Error converting str to int", err)
				msg.Text = "Number of orders was not a valid integer. Use /generate {integer}"
				continue
			}
			generateResponse, err := callOrderGenerationServiceBulk(ordersCount)
			if err != nil {
				log.Println("Error generating order:", err)
				continue
			}
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, generateResponse)
		case "status":
			msg.Text = "Server status: "
		default:
			msg.Text = "Unknown command"
		}

		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}
}

func callOrderGenerationService() (string, error) {
	response, err := http.Get("http://order-generation-service:8081//fetch-products")
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

func callOrderGenerationServiceBulk(i int) (string, error) {
	response, err := http.Get("http://order-generation-service:8081/generate-orders/" + strconv.Itoa(i))
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
