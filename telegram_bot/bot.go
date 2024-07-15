package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
	"os"
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
			msg.Text = "I understand /generate and /status."
		case "generate":
			msg.Text = "executing generate orders..."
			// call order-generation-service
			orders, err := callOrderGenerationService()
			if err != nil {
				log.Println("Error generating order:", err)
				continue
			}
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, orders)
		case "gene5":
			msg.Text = "executing generate orders..."
			// call order-generation-service
			orders, err := callOrderGenerationServiceBulk()
			if err != nil {
				log.Println("Error generating order:", err)
				continue
			}
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, orders)
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
	response, err := http.Get("http://order-generation-service:8081/generate-order")
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
func callOrderGenerationServiceBulk() (string, error) {
	response, err := http.Get("http://order-generation-service:8081/generate-orders/5")
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
