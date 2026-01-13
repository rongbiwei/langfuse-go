package observer

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

type command int

const (
	commanFlush command = iota
	commandFlushAndWait
)

const (
	defaultTickerPeriod = 1 * time.Second
	maxHandleGoroutines = 5 // 最大并发处理协程数
)

type handler[T any] struct {
	queue        *queue[T]
	fn           EventHandler[T]
	commandCh    chan command
	tickerPeriod time.Duration
	semaphore    chan struct{}  // 协程信号量
	wg           sync.WaitGroup // 等待所有 handle goroutine 完成
	closed       atomic.Bool    // 标记是否已关闭
}

func newHandler[T any](queue *queue[T], fn EventHandler[T]) *handler[T] {
	return &handler[T]{
		queue:        queue,
		fn:           fn,
		commandCh:    make(chan command),
		tickerPeriod: defaultTickerPeriod,
		semaphore:    make(chan struct{}, maxHandleGoroutines),
	}
}

func (h *handler[T]) withTick(period time.Duration) *handler[T] {
	h.tickerPeriod = period
	return h
}

func (h *handler[T]) listen(ctx context.Context) {
	ticker := time.NewTicker(h.tickerPeriod)
	defer ticker.Stop()
	defer h.markClosed()

	for {
		select {
		case <-ctx.Done():
			// 等待所有正在执行的 handle 完成
			h.wg.Wait()
			// 处理队列中剩余的数据
			h.handle(ctx)
			return
		case <-ticker.C:
			// 尝试获取信号量，获取不到则跳过本次
			select {
			case h.semaphore <- struct{}{}:
				h.wg.Add(1)
				go func() {
					defer func() {
						<-h.semaphore
						h.wg.Done()
					}()
					h.handle(ctx)
				}()
			default:
				// 已达到最大并发数，跳过本次
			}
		case cmd, ok := <-h.commandCh:
			if !ok {
				return
			}

			if cmd == commanFlush {
				h.handle(ctx)
			} else if cmd == commandFlushAndWait {
				// 等待所有正在执行的 handle 完成
				h.wg.Wait()
				// 最后再处理一次队列中剩余的数据
				h.handle(ctx)
				close(h.commandCh)
			}
		}
	}
}

func (h *handler[T]) markClosed() {
	h.closed.Store(true)
}

func (h *handler[T]) handle(ctx context.Context) {
	events := h.queue.All()
	if len(events) == 0 {
		return
	}
	res := h.fn(ctx, events)
	if len(res) != 0 {
		for _, e := range res {
			// 使用非阻塞方式入队，队列满时丢弃失败事件
			h.queue.TryEnqueue(e)
		}
	}
}

func (h *handler[T]) flush() {
	if h.closed.Load() {
		return
	}
	select {
	case h.commandCh <- commanFlush:
	default:
		// channel 已关闭或满，忽略
	}
}

func (h *handler[T]) flushAndWait() {
	if h.closed.Load() {
		return
	}
	
	// 使用 recover 防止向已关闭 channel 发送导致 panic
	defer func() {
		recover()
	}()
	
	h.commandCh <- commandFlushAndWait
	// 等待 channel 关闭
	for range h.commandCh {
	}
}
