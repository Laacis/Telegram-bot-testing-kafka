package telegram_bot

import (
	"log"
	"os"
	"strconv"
	"testing"
)

func init() {
	//setting up
	os.Setenv("TELEGRAM_BOT_TOKEN", "test_token")
	os.Setenv("BOT_TIMEOUT", strconv.Itoa(20))
	os.Setenv("DEBUG_BOT", "true")
	os.Setenv("PRODUCER_MANAGER_SERVICE_PREFIX", "localhost:8082")
	os.Setenv("ORDER_SERVICE_PREFIX", "localhost:9192")
	InitConfig()
}

func TestBotToken(t *testing.T) {
	expectedToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	got := BotToken()
	if got != expectedToken {
		t.Errorf("BotToken() = %v; expected  %v", got, expectedToken)
	}
}

func TestBotUpdateConfig(t *testing.T) {
	updateConfig := BotUpdateConfig()
	expectedTimeout, err := strconv.Atoi(os.Getenv("BOT_TIMEOUT"))
	if err != nil {
		log.Printf("Invalid BOT_TIMEOUT value: %v", err)
	}
	returnedTimeout := updateConfig.Timeout
	if returnedTimeout != expectedTimeout {
		t.Errorf("TestBotUpdateConfig() returned Timeout: %v; expected %v", returnedTimeout, expectedTimeout)
	}
}

func TestGetEndpoint(t *testing.T) {
	kafkaManagerEndpointPrefix := os.Getenv("PRODUCER_MANAGER_SERVICE_PREFIX")
	orderServiceEndpointPrefix := os.Getenv("ORDER_SERVICE_PREFIX")
	tests := []struct {
		command   string
		params    []int
		expected  string
		shouldErr bool
	}{
		{"producerUp", nil, kafkaManagerEndpointPrefix + "/producer/up", false},
		{"producerDown", nil, kafkaManagerEndpointPrefix + "/producer/down", false},
		{"producerStatus", nil, kafkaManagerEndpointPrefix + "/producer/status", false},
		{"sendAll", nil, orderServiceEndpointPrefix + "/orders/send/all", false},
		{"send", []int{99}, orderServiceEndpointPrefix + "/orders/send/99", false},
		{"generate", []int{56}, orderServiceEndpointPrefix + "/generate-orders/56", false},
		{"send", nil, "", true},
		{"generate", nil, "", true},
		{"unknown", nil, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.command, func(t *testing.T) {
			got, err := GetEndpoint(tt.command, tt.params...)
			if (err != nil) != tt.shouldErr {
				t.Errorf("GetEndpoint() error = %v, shouldErr %v", err, tt.shouldErr)
				return
			}
			if got != tt.expected {
				t.Errorf("GetEndpoint() = %v, expected %v", got, tt.expected)
			}
		})
	}
}
