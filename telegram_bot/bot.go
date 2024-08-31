package main

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"strconv"
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

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		commandString, arg, err := extractArgs(update)
		if err != nil {
			log.Printf("Error extracting args: %v", err)
			msg.Text = "An error occurred while handling your request"
		} else {
			cmd, err := commands.Create(commandString, arg, cfg)
			if err != nil {
				log.Printf("Error crafting cmd: %v", err)
			}
			if cmd != nil {
				response, err := cmd.Execute(httpClient)
				if err != nil {
					log.Printf("Error handling response: %v", err)
					msg.Text = "An error occurred while handling your request"
				} else {
					msg.Text = responseMessage(err, response)
				}
			} else {
				str := fmt.Sprintf("Unknown cmd %s! I understand:\n /producerUp \n/producerDown \n/producerStatus \n/generate X \n/send X.", commandString)
				msg.Text = str
			}
		}

		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}
}

func responseMessage(err error, response []byte) string {
	var data interface{}
	var message string
	err = json.Unmarshal(response, &data)
	if err != nil {
		log.Fatalf("error unmarshalling JSON: %v", err)
	}
	if d, ok := data.(map[string]interface{}); ok {
		if m, ok := d["Message"].(string); ok {
			message = m
		}
	}
	return message
}

func extractArgs(update tgbotapi.Update) (string, int, error) {
	commandString := update.Message.Command()
	args := strings.Fields(update.Message.Text)[1:]
	if len(args) > 0 {
		arg, err := strconv.Atoi(args[0])
		if err != nil {
			log.Printf("Error converting arg to int: %v", err)
			return "", 0, err
		}
		return commandString, arg, nil
	}
	return commandString, 0, nil
}
