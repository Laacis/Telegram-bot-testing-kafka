package telegram_bot

import (
	"fmt"
	"io"
	"net/http"
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

func Create(command string, arg int, endpointGetter EndpointGetter) (*Command, error) {
	var endpoint string
	var err error
	if arg == 0 {
		endpoint, err = endpointGetter.GetEndpoint(command)
	} else {
		endpoint, err = endpointGetter.GetEndpoint(command, arg)
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
	response, err := c.Get(h.endpoint)
	if err != nil {
		return "", fmt.Errorf("failed to perform GET request to %s: %w", h.endpoint, err)
	}
	defer func() { _ = response.Body.Close() }()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return "", fmt.Errorf("HTTP request failed with status code %d", response.StatusCode)
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body from %s: %w", h.endpoint, err)
	}

	return string(data), nil
}
