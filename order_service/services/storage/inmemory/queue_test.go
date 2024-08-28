package order_generation_service

import (
	"testing"
)

func TestNewQueue(t *testing.T) {
	cap := 1
	expected := Queue[string]{
		data: make([]string, 0),
		cap:  cap,
	}
	got := NewQueue[string](cap)

	if got.cap != expected.cap {
		t.Errorf("NewQueue() returned with cap: %d, expected: %d", got.cap, expected.cap)
	}
}

func TestQueue_Enqueue(t *testing.T) {
	expectedStr := "test"
	cap := 5
	expectedBool := true
	queue := &Queue[string]{
		data: make([]string, cap),
		cap:  cap,
	}
	result := queue.Enqueue(expectedStr)
	if !result {
		t.Errorf("Enqueue() returned %v, expected: %v", result, expectedBool)
	}
	if queue.size != 1 {
		t.Errorf("Enqueue() size: %v, erxepted: %v", queue.size, 1)
	}
	if queue.data[0] != expectedStr {
		t.Errorf("Enqueue() returned data[0]: %s, expected %s", queue.data[0], expectedStr)
	}
}

func TestQueue_IsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		queue    *Queue[string]
		expected bool
	}{
		{"QueueIsEmpty", &Queue[string]{data: make([]string, 1), cap: 1}, true},
		{"QueueIsNotEmpty", &Queue[string]{data: append(make([]string, 1), "test"), cap: 1, size: 1}, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.queue.IsEmpty()

			if got != test.expected {
				t.Errorf("IsEmpty() returned %v, expected to be %v", got, test.expected)
			}
		})
	}
}

func TestQueue_IsFull(t *testing.T) {
	tests := []struct {
		name     string
		queue    *Queue[string]
		expected bool
	}{
		{"QueueIsFull", &Queue[string]{data: append(make([]string, 1), "test"), cap: 1, size: 1}, true},
		{"QueueIsNotFull", &Queue[string]{data: append(make([]string, 1, 2), "test"), cap: 2, size: 1}, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			returned := test.queue.IsFull()

			if returned != test.expected {
				t.Errorf("IsFull() returned: %v, expected %v", returned, test.expected)
			}
		})
	}
}
