package order_generation_service

import (
	"reflect"
	"testing"
)

func TestNewInMemoryStorage(t *testing.T) {
	expected := &InMemoryStorage[string]{items: make([]string, 0)}
	got := NewInMemoryStorage[string]()

	if len(got.items) != 0 {
		t.Errorf("NewInMemoryStorage() = %v, expected %v", got.items, expected.items)
	}
	expectedType := reflect.TypeOf(expected.items)
	gotType := reflect.TypeOf(got.items)
	if expectedType != gotType {
		t.Errorf("NewInMemoryStorage() type is %v, expected %v", gotType, expectedType)
	}
}

func TestInMemoryStorage_Add(t *testing.T) {
	expectedString := "test"
	storage := &InMemoryStorage[string]{items: make([]string, 0)}
	storage.Add(expectedString)

	storageItems := len(storage.items)
	if storageItems != 1 {
		t.Errorf("Add(), expected to have 1 item in storage, got: %d", storageItems)
	}

	firstItem := storage.items[0]
	if expectedString != firstItem {
		t.Errorf("Expected item %v, but got %v", expectedString, firstItem)
	}
}

func TestInMemoryStorage_AllRecords(t *testing.T) {
	testString := "one"
	storage := &InMemoryStorage[string]{items: append(make([]string, 0), testString)}

	got := storage.AllRecords()
	//got is *[]string
	if (*got)[0] != testString {
		t.Errorf("AllReccords() first item %s, is not as expected: %s", (*got)[0], testString)
	}
}

func TestInMemoryStorage_Length(t *testing.T) {
	tests := []struct {
		name        string
		items       []string
		lenExpected int
	}{
		{"empty", make([]string, 0), 0},
		{"notEmpty", append(make([]string, 0), "one", "two"), 2},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			storage := &InMemoryStorage[string]{items: test.items}
			got := storage.Length()
			if got != test.lenExpected {
				t.Errorf("Length() returned %d, expected: %d", got, test.lenExpected)
			}
		})
	}
}
