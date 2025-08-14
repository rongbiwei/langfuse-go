package observer

import "sync"

type queue[T any] struct {
	sync.Mutex
	items []T
}

// Enqueue 入队
func (q *queue[T]) Enqueue(item T) {
	q.Lock()
	defer q.Unlock()
	q.items = append(q.items, item)
}

// Dequeue 出队
func (q *queue[T]) Dequeue() T {
	q.Lock()
	defer q.Unlock()
	if len(q.items) == 0 {
		var zero T
		return zero
	}
	item := q.items[0]
	q.items = q.items[1:]
	return item
}

// Len 长度
func (q *queue[T]) Len() int {
	q.Lock()
	defer q.Unlock()
	return len(q.items)
}

func newQueue[T any]() *queue[T] {
	return &queue[T]{}
}

// Clear 清空
func (q *queue[T]) Clear() {
	q.Lock()
	defer q.Unlock()
	q.items = []T{}
}

// All 获取所有
func (q *queue[T]) All() []T {
	q.Lock()
	defer q.Unlock()
	items := q.items
	q.items = []T{}
	return items
}
