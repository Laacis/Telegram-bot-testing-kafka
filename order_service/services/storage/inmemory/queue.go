package order_generation_service

type Queue[T any] struct {
	data []T
	head int
	tail int
	size int
	cap  int
}

func NewQueue[T any](cap int) *Queue[T] {
	return &Queue[T]{
		data: make([]T, 0),
		cap:  cap,
	}
}

func (q *Queue[T]) Enqueue(value T) bool {
	if q.size == q.cap {
		return false
	}

	q.data[q.tail] = value
	q.tail = (q.tail + 1) % q.cap
	q.size++
	return true
}

func (q *Queue[T]) Dequeue() (T, bool) {
	var zero T
	if q.size == 0 {
		return zero, false
	}

	value := q.data[q.head]
	q.head = (q.head + 1) % q.cap
	q.size--
	return value, false
}

func (q *Queue[T]) IsFull() bool {
	return q.size == q.cap
}

func (q *Queue[T]) IsEmpty() bool {
	return q.size == 0
}
