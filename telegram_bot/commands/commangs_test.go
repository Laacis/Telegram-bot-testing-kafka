package telegram_bot

import (
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"
)

type MockEndpointGetter struct{}
type MockHttpClient struct{}

func (m *MockEndpointGetter) GetEndpoint(command string, args ...int) (string, error) {
	if command == "fail" {
		return "", errors.New("mock endpoint failure")
	}
	var str string
	if len(args) > 0 {
		str = command + "/endpoint/" + strconv.Itoa(args[0])
	} else {
		str = command + "/endpoint/"
	}
	return str, nil
}

func (client *MockHttpClient) Get(endpoint string) (*http.Response, error) {
	if endpoint == "error" {
		return nil, errors.New("mocked failure")
	}
	statusCode, err := strconv.Atoi(endpoint)
	if err != nil {
		return nil, errors.New("error converting mocked status code")
	}
	response := &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(strings.NewReader(endpoint)),
	}
	return response, nil
}

func TestCreateCommand(t *testing.T) {
	tests := []struct {
		caseName   string
		commandStr string
		arg        int
		expected   *Command
		shouldErr  bool
	}{
		{"successfullyCreateCommandNoArgs", "test", 0, &Command{endpoint: "test/endpoint/"}, false},
		{"successfullyCreateCommandOneArg", "test", 1, &Command{endpoint: "test/endpoint/1"}, false},
		{"failCreateCommandGetterError", "fail", 1, nil, true},
	}

	for _, test := range tests {
		t.Run(test.caseName, func(t *testing.T) {
			mockEndpointGetter := new(MockEndpointGetter)

			returned, err := Create(test.commandStr, test.arg, mockEndpointGetter)
			if (err != nil) != test.shouldErr {
				t.Errorf("CreateCommand() error = %v, shouldErr %v", err, test.shouldErr)
				return
			}
			if !endpointComparor(returned, test.expected) {
				t.Errorf("CraftCommand() returned endpoint doesn't match the expected.")
			}
		})
	}
}

func endpointComparor(returned *Command, expected *Command) bool {
	if returned == nil && expected == nil {
		return true
	}
	return returned.endpoint == expected.endpoint
}

func TestCommand_Execute(t *testing.T) {
	tests := []struct {
		caseName  string
		command   *Command
		expected  string
		shouldErr bool
	}{
		{"successfullyExecute", &Command{endpoint: "200"}, "200", false},
		{"failExecuteClientErr", &Command{endpoint: "error"}, "", true},
		{"failExecuteResponseCode404", &Command{endpoint: "404"}, "", true},
		{"failExecuteResponseCode500", &Command{endpoint: "500"}, "", true},
	}

	for _, test := range tests {
		t.Run(test.caseName, func(t *testing.T) {
			mockHttpClient := new(MockHttpClient)
			response, err := test.command.Execute(mockHttpClient)
			if (err != nil) != test.shouldErr {
				t.Errorf("CreateCommand() error = %v, shouldErr %v", err, test.shouldErr)
				return
			}

			if test.expected != response {
				t.Errorf("command.Excute() returned %v, but was expected %v.", response, test.expected)
			}

		})
	}
}
