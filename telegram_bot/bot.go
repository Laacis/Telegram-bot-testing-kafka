package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"strings"
	commands "telegram_bot/commands"
	config "telegram_bot/config"
)

func init() {
	loadEnv()
	config.InitConfig()
}

func main() {
	bot := botSetup(config.BotToken())
	updates := bot.GetUpdatesChan(config.BotUpdateConfig())

	cfg := new(config.ConfigData)
	httpClient := new(http.Client)

	updateHandler(updates, bot, cfg, httpClient)
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

func updateHandler(updates tgbotapi.UpdatesChannel, bot *tgbotapi.BotAPI, cfg *config.ConfigData, httpClient commands.HTTPClient) {
	for update := range updates {
		if update.Message == nil || !update.Message.IsCommand() {
			continue
		}

		// Craft command-string
		commandArgs := craftCommandString(update)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		// Get command Object
		commandObj, err := commands.CraftCommand(commandArgs, cfg)
		if err != nil {
			log.Printf("Error handling handler: %v", err)
		}
		// TODO verify command exists, has endpoint or is internal
		// Execute the command
		response, err := commandObj.Execute(httpClient)
		if err != nil {
			log.Printf("Error handling handler: %v", err)
			response = "An error occurred while handling your request" + err.Error()
		}
		msg.Text = response
		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}
}

func craftCommandString(update tgbotapi.Update) []string {
	command := update.Message.Command()
	var args []string
	x := strings.Fields(update.Message.Text)[1:]
	args = append(args, command)
	args = append(args, x...)
	return args
}
