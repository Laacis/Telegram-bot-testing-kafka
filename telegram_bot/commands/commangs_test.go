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
	if endpoint != "localhost" {
		return &http.Response{
			StatusCode: 404,
			Body:       io.NopCloser(strings.NewReader("HostNotFound")),
		}, errors.New("mocked failure 404")
	}
	mockedResponseStr := "Executed successfully"
	response := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader(mockedResponseStr)),
	}
	return response, nil
}

func TestCreateCommand(t *testing.T) {
	tests := []struct {
		caseName  string
		command   string
		params    string
		expected  *Command
		shouldErr bool
	}{
		{"successfullyCreateCommand", "test", "1", &Command{endpoint: "test/endpoint/1"}, false},
		{"failCreateCommandWithNotIntArg", "test", "a", nil, true},
		{"failCreateCommandGetterError", "fail", "1", nil, true},
	}

	for _, test := range tests {
		t.Run(test.caseName, func(t *testing.T) {
			mockEndpointGetter := new(MockEndpointGetter)
			args := append(make([]string, 0, 2), test.command, test.params)
			returned, err := CreateCommand(args, mockEndpointGetter)
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
		{"successfullyExecute", &Command{endpoint: "localhost"}, "Executed successfully", false},
		{"failExecuteClientErr", &Command{endpoint: "failHost"}, "", true},
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
