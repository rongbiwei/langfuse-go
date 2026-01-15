package observer

import (
	"encoding/json"
	"sync"
)

const (
	// DefaultMaxMemoryBytes 默认最大内存限制 128MB
	DefaultMaxMemoryBytes = 50 * 1024 * 1024
)

type queue[T any] struct {
	mu          sync.Mutex
	cond        *sync.Cond
	items       []T
	currentSize int64 // 当前队列数据大小（字节）
	maxSize     int64 // 最大允许大小（字节）
}

// Enqueue 入队，当队列大小超过限制时阻塞
func (q *queue[T]) Enqueue(item T) {
	itemSize := int64(q.estimateItemSize(item))

	q.mu.Lock()
	defer q.mu.Unlock()

	// 当队列大小超过限制时，等待直到有空间
	for q.currentSize+itemSize > q.maxSize && len(q.items) > 0 {
		q.cond.Wait()
	}

	q.items = append(q.items, item)
	q.currentSize += itemSize
}

// TryEnqueue 尝试入队，如果队列已满则返回 false（不阻塞）
func (q *queue[T]) TryEnqueue(item T) bool {
	itemSize := int64(q.estimateItemSize(item))

	q.mu.Lock()
	defer q.mu.Unlock()

	// 队列已满，直接返回失败
	if q.currentSize+itemSize > q.maxSize && len(q.items) > 0 {
		return false
	}

	q.items = append(q.items, item)
	q.currentSize += itemSize
	return true
}

// estimateItemSize 使用 json.Marshal 估算元素大小
func (q *queue[T]) estimateItemSize(item T) int {
	data, err := json.Marshal(item)
	if err != nil {
		return 1024 // 序列化失败时返回保守估算值
	}
	return len(data)
}

// Dequeue 出队
func (q *queue[T]) Dequeue() T {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.items) == 0 {
		var zero T
		return zero
	}
	item := q.items[0]
	q.items = q.items[1:]
	q.currentSize -= int64(q.estimateItemSize(item))
	if q.currentSize < 0 {
		q.currentSize = 0
	}
	q.cond.Signal()
	return item
}

// Len 长度
func (q *queue[T]) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.items)
}

// Size 返回当前队列数据大小（字节）
func (q *queue[T]) Size() int64 {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.currentSize
}

func newQueue[T any]() *queue[T] {
	q := &queue[T]{
		maxSize: DefaultMaxMemoryBytes,
	}
	q.cond = sync.NewCond(&q.mu)
	return q
}

// newQueueWithMaxSize 创建指定最大内存限制的队列
func newQueueWithMaxSize[T any](maxSize int64) *queue[T] {
	q := &queue[T]{
		maxSize: maxSize,
	}
	q.cond = sync.NewCond(&q.mu)
	return q
}

// Clear 清空
func (q *queue[T]) Clear() {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.items = []T{}
	q.currentSize = 0
	q.cond.Broadcast()
}

// All 获取所有
func (q *queue[T]) All() []T {
	q.mu.Lock()
	defer q.mu.Unlock()
	items := q.items
	q.items = []T{}
	q.currentSize = 0
	q.cond.Broadcast()
	return items
}
