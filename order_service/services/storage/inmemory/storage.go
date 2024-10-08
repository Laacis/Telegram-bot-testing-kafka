package order_generation_service

type InMemoryStorageInterface[T any] interface {
	Add(item T)
	Length() int
	AllRecords() []T
}

type InMemoryStorage[T any] struct {
	items []T
}

func NewInMemoryStorage[T any]() *InMemoryStorage[T] {
	return &InMemoryStorage[T]{
		items: make([]T, 0),
	}
}

func (s *InMemoryStorage[T]) Add(item T) {
	s.items = append(s.items, item)
}

func (s *InMemoryStorage[T]) AllRecords() *[]T {
	return &s.items
}

func (s *InMemoryStorage[T]) Length() int {
	return len(s.items)
}

type myInterface[T any] interface {
	Next() (T, bool)
}

func (s *InMemoryStorage[T]) Next() (*T, bool) {
	if len(s.items) == 0 {
		return nil, false
	}
	next := s.items[0]
	s.items = s.items[1:]
	return &next, true
}
