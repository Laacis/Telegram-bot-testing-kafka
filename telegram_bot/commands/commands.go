package telegram_bot

import (
	"io"
	"log"
	"net/http"
	"strconv"
)

type Command struct {
	endpoint string
}

type EndpointGetter interface {
	GetEndpoint(command string, args ...int) (string, error)
}

type HTTPClient interface {
	Get(url string) (*http.Response, error)
}

func handleHelp(args []string) (string, error) {
	return "I understand:\n /producerUp \n/producerDown \n/producerStatus \n/generate X \n/send X.", nil
}

func CraftCommand(args []string, endpointGetter EndpointGetter) (*Command, error) {
	//command comes as first args[0]
	command := args[0]
	//possible int args following
	var intArgs int
	var err error
	if len(args) > 1 {
		intArgs, err = strconv.Atoi(args[1])
		if err != nil {
			log.Printf("Error converting arg to int: %v", err)
			return nil, err
		}
	}
	var endpoint string

	if intArgs == 0 {
		endpoint, err = endpointGetter.GetEndpoint(command)
	} else {
		endpoint, err = endpointGetter.GetEndpoint(command, intArgs)
	}
	if err != nil {
		return nil, err
	}

	result := Command{
		endpoint: endpoint,
	}
	return &result, nil
}

func (h *Command) Execute(c HTTPClient) (string, error) {
	resultStr := ""
	response, err := c.Get(h.endpoint)
	if err != nil {
		return "", err
	}
	defer func() { _ = response.Body.Close() }()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	resultStr = string(data)
	return resultStr, nil
}
