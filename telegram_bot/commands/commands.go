package telegram_bot

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	config "telegram_bot/config"
)

var CommandHandlers = map[string]func([]string) (string, error){
	"help":           handleHelp,
	"producerUp":     simpleGetHandlerNoArguments,
	"producerDown":   simpleGetHandlerNoArguments,
	"producerStatus": simpleGetHandlerNoArguments,
	"sendAll":        simpleGetHandlerNoArguments,
	"send":           simpleGetHandlerOneArgument,
	"generate":       simpleGetHandlerOneArgument,
	"status":         handleStatus,
}

func handleHelp(args []string) (string, error) {
	return "I understand:\n /producerUp \n/producerDown \n/producerStatus \n/generate X \n/send X.", nil
}

func simpleGetHandlerNoArguments(args []string) (string, error) {
	command := strings.Split(args[0], " ")[0]
	endpoint, err := config.GetEndpoint(command)
	if err != nil {
		return "", err
	}

	response, err := http.Get(endpoint)
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

func simpleGetHandlerOneArgument(args []string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("missing parameter for generate command")
	}

	numOrders, err := strconv.Atoi(args[0])
	if err != nil {
		return "", fmt.Errorf("invalid number of orders: %s", args[0])
	}

	endpoint, err := config.GetEndpoint("generate", numOrders)
	if err != nil {
		return "", err
	}

	response, err := http.Get(endpoint)
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

func handleStatus(args []string) (string, error) {
	return "Server status: TODO", nil
}
