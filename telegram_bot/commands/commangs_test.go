package telegram_bot

import (
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"testing"
)

type MockEndpointGetter struct{}
type MockHttpClient struct{}

func (m *MockEndpointGetter) GetEndpoint(command string, args ...int) (string, error) {
	var str string
	if len(args) > 0 {
		str = command + "/endpoint/" + strconv.Itoa(args[0])
	} else {
		str = command + "/endpoint/"
	}
	return str, nil
}

func (client *MockHttpClient) Get(endpoint string) (*http.Response, error) {
	mockedResponseStr := "Executed successfully"
	response := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader(mockedResponseStr)),
	}
	return response, nil
}

func TestCreateCommand(t *testing.T) {
	mockEndpointGetter := new(MockEndpointGetter)
	command := "producerUpt"
	intArgs := "1"
	args := append(make([]string, 0, 2), command, intArgs)
	expectedStr := command + "/endpoint/" + intArgs
	expectedObject := Command{
		endpoint: expectedStr,
	}
	returned, err := CreateCommand(args, mockEndpointGetter)
	if err != nil {
		log.Printf("Error while testing : %v", err)
	}
	if returned.endpoint != expectedObject.endpoint {
		t.Errorf("CraftCommand() returned %v, but was expected %v.", returned.endpoint, expectedObject.endpoint)
	}
}

func TestCommand_Execute(t *testing.T) {
	cmd := Command{
		endpoint: "localhost",
	}
	expectedStr := "Executed successfully"
	mockHttpClient := new(MockHttpClient)
	response, err := cmd.Execute(mockHttpClient)
	if err != nil {
		log.Printf("Error while testing : %v", err)
	}
	if expectedStr != response {
		t.Errorf("command.Excute() returned %v, but was expected %v.", response, expectedStr)
	}

}
