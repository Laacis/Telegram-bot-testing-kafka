package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"log"
	"strings"
	commands "telegram_bot/commands"
	config "telegram_bot/config"
)

func main() {
	loadEnv()
	bot := botSetup(config.BotToken())
	updates := bot.GetUpdatesChan(config.BotUpdateConfig())

	updateHandler(updates, bot)
}

func botSetup(botToken string) *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = config.DebugBot
	log.Printf("Authorized on account %s", bot.Self.UserName)
	return bot
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func updateHandler(updates tgbotapi.UpdatesChannel, bot *tgbotapi.BotAPI) {
	for update := range updates {
		if update.Message == nil || !update.Message.IsCommand() {
			continue
		}

		command := update.Message.Command()
		var args []string
		x := strings.Fields(update.Message.Text)[1:]
		args = append(args, command)
		args = append(args, x...)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		handler, exists := commands.CommandHandlers[command]
		if exists {
			response, err := handler(args)
			if err != nil {
				log.Printf("Error handling command: %s: %v", command, err)
				response = "An error occurred while handling your request" + err.Error()
			}
			msg.Text = response
		} else {
			msg.Text = "Unknown command"
		}
		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}
}
