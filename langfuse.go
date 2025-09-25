package langfuse

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/rongbiwei/langfuse-go/internal/pkg/api"
	"github.com/rongbiwei/langfuse-go/internal/pkg/observer"
	"github.com/rongbiwei/langfuse-go/model"
)

const (
	defaultFlushInterval = 500 * time.Millisecond
	retryInterval        = 3
	// batchSize 每次批量发送的数据量
	batchSize = 3 * 1024 * 1024
)

// Langfuse 跟踪对象
type Langfuse struct {
	flushInterval time.Duration
	client        *api.Client
	observer      *observer.Observer[model.IngestionEvent]
	location      *time.Location
}

// New 创建一个新的Langfuse
func New(ctx context.Context, parallel int) *Langfuse {
	client := api.New()
	loc, _ := time.LoadLocation("Asia/Shanghai")
	l := &Langfuse{
		flushInterval: defaultFlushInterval,
		client:        client,
		location:      loc,
		observer: observer.NewObserver(
			ctx,
			func(ctx context.Context, events []model.IngestionEvent) []model.IngestionEvent {
				failEvents := make([]model.IngestionEvent, 0)
				if len(events) == 0 {
					return nil
				}
				pushDataBatch(ctx, parallel, client, events, failEvents)
				return failEvents
			},
		),
	}
	return l
}

// pushData 推送数据
func pushData(ctx context.Context, client *api.Client, index int, events, failEvents []model.IngestionEvent) {
	if index >= len(events) || len(events) == 0 {
		return
	}
	batchData := make([]model.IngestionEvent, 0)
	currentSize := 0
	for i := index; i < len(events); i++ {
		byteArr, _ := json.Marshal(events[i])
		if currentSize+len(byteArr) > batchSize {
			if i == index {
				// 單條數據超過大小
				batchData = append(batchData, events[i])
				index++
			}
			break
		}
		currentSize += len(byteArr)
		batchData = append(batchData, events[i])
		index++
	}
	if err := ingest(ctx, client, batchData); err != nil {
		log.Println("ingest error:" + err.Error())
		for _, e := range batchData {
			if e.FailCount < 3 {
				e.FailCount++
				failEvents = append(failEvents, e)
			}
		}
	}
	pushData(ctx, client, index, events, failEvents)
}

// pushDataBatch 推送数据--- 批量
func pushDataBatch(ctx context.Context, parallel int, client *api.Client, events, failEvents []model.IngestionEvent) {
	if parallel <= 0 {
		parallel = 5
	}
	var wg sync.WaitGroup
	maxGoroutines, index := parallel, 0
	goroutineSemaphore := make(chan struct{}, maxGoroutines)

	for index < len(events) && len(events) > 0 {
		batchData := make([]model.IngestionEvent, 0)
		currentSize := 0
		for i := index; i < len(events); i++ {
			byteArr, _ := json.Marshal(events[i])
			if currentSize+len(byteArr) > batchSize {
				if i == index {
					// 單條數據超過大小
					batchData = append(batchData, events[i])
					index++
				}
				break
			}
			currentSize += len(byteArr)
			batchData = append(batchData, events[i])
			index++
		}

		if len(batchData) > 0 {
			goroutineSemaphore <- struct{}{}
			wg.Add(1)
			go func(batch []model.IngestionEvent) {
				defer func() {
					<-goroutineSemaphore
					wg.Done()
				}()
				if err := ingest(ctx, client, batch); err != nil {
					log.Println("ingest error:" + err.Error())
					//for _, e := range batch {
					//	if e.FailCount < 3 {
					//		e.FailCount++
					//		failEvents = append(failEvents, e)
					//	}
					//}
				}
			}(batchData)
		}
	}
	wg.Wait()
}

func (l *Langfuse) WithFlushInterval(d time.Duration) *Langfuse {
	l.flushInterval = d
	return l
}

func ingest(ctx context.Context, client *api.Client, events []model.IngestionEvent) error {
	req := api.Ingestion{
		Batch: events,
	}

	res := api.IngestionResponse{}
	return client.Ingestion(ctx, &req, &res)
}

// Trace 构建跟踪
func (l *Langfuse) Trace(t *model.Trace) (*model.Trace, error) {
	t.ID = buildID(&t.ID)
	now := time.Now()
	if l.location != nil {
		now = now.In(l.location)
	}
	l.observer.Dispatch(
		model.IngestionEvent{
			ID:        buildID(nil),
			Type:      model.IngestionEventTypeTraceCreate,
			Timestamp: now,
			Body:      t,
		},
	)
	return t, nil
}

// TraceWithTime 构建跟踪并指定时间戳
func (l *Langfuse) TraceWithTime(t *model.Trace, timestamp time.Time) (*model.Trace, error) {
	t.ID = buildID(&t.ID)
	l.observer.Dispatch(
		model.IngestionEvent{
			ID:        t.ID,
			Type:      model.IngestionEventTypeTraceCreate,
			Timestamp: timestamp,
			Body:      t,
		},
	)
	return t, nil
}

// Generation 构建生成
func (l *Langfuse) Generation(g *model.Generation, parentID *string) (*model.Generation, error) {
	if g.TraceID == "" {
		traceID, err := l.createTrace(g.Name)
		if err != nil {
			return nil, err
		}

		g.TraceID = traceID
	}

	g.ID = buildID(&g.ID)

	if parentID != nil {
		g.ParentObservationID = *parentID
	}
	now := time.Now()
	if l.location != nil {
		now = now.In(l.location)
	}
	l.observer.Dispatch(
		model.IngestionEvent{
			ID:        buildID(nil),
			Type:      model.IngestionEventTypeGenerationCreate,
			Timestamp: now,
			Body:      g,
		},
	)
	return g, nil
}

// GenerationWithTime 构建生成并指定时间戳
func (l *Langfuse) GenerationWithTime(g *model.Generation, parentID *string, timestamp time.Time) (*model.Generation, error) {
	if g.TraceID == "" {
		traceID, err := l.createTrace(g.Name)
		if err != nil {
			return nil, err
		}

		g.TraceID = traceID
	}

	g.ID = buildID(&g.ID)

	if parentID != nil {
		g.ParentObservationID = *parentID
	}

	l.observer.Dispatch(
		model.IngestionEvent{
			ID:        g.ID,
			Type:      model.IngestionEventTypeGenerationCreate,
			Timestamp: timestamp,
			Body:      g,
		},
	)
	return g, nil
}

// GenerationEnd 结束一个生成
func (l *Langfuse) GenerationEnd(g *model.Generation) (*model.Generation, error) {
	if g.ID == "" {
		return nil, fmt.Errorf("generation ID is required")
	}

	if g.TraceID == "" {
		return nil, fmt.Errorf("trace ID is required")
	}
	now := time.Now()
	if l.location != nil {
		now = now.In(l.location)
	}
	l.observer.Dispatch(
		model.IngestionEvent{
			ID:        buildID(nil),
			Type:      model.IngestionEventTypeGenerationUpdate,
			Timestamp: now,
			Body:      g,
		},
	)

	return g, nil
}

// GenerationEndWithTime 结束一个生成并指定时间戳
func (l *Langfuse) GenerationEndWithTime(g *model.Generation, timestamp time.Time) (*model.Generation, error) {
	if g.ID == "" {
		return nil, fmt.Errorf("generation ID is required")
	}

	if g.TraceID == "" {
		return nil, fmt.Errorf("trace ID is required")
	}

	l.observer.Dispatch(
		model.IngestionEvent{
			ID:        g.ID,
			Type:      model.IngestionEventTypeGenerationUpdate,
			Timestamp: timestamp,
			Body:      g,
		},
	)
	return g, nil
}

// Score 构建分数
func (l *Langfuse) Score(s *model.Score) (*model.Score, error) {
	if s.TraceID == "" {
		return nil, fmt.Errorf("trace ID is required")
	}
	s.ID = buildID(&s.ID)
	now := time.Now()
	if l.location != nil {
		now = now.In(l.location)
	}
	l.observer.Dispatch(
		model.IngestionEvent{
			ID:        buildID(nil),
			Type:      model.IngestionEventTypeScoreCreate,
			Timestamp: now,
			Body:      s,
		},
	)
	return s, nil
}

// Span 构建span
func (l *Langfuse) Span(s *model.Span, parentID *string) (*model.Span, error) {
	if s.TraceID == "" {
		traceID, err := l.createTrace(s.Name)
		if err != nil {
			return nil, err
		}

		s.TraceID = traceID
	}

	s.ID = buildID(&s.ID)

	if parentID != nil {
		s.ParentObservationID = *parentID
	}
	now := time.Now()
	if l.location != nil {
		now = now.In(l.location)
	}
	l.observer.Dispatch(
		model.IngestionEvent{
			ID:        buildID(nil),
			Type:      model.IngestionEventTypeSpanCreate,
			Timestamp: now,
			Body:      s,
		},
	)

	return s, nil
}

// SpanEnd 结束一个span
func (l *Langfuse) SpanEnd(s *model.Span) (*model.Span, error) {
	if s.ID == "" {
		return nil, fmt.Errorf("generation ID is required")
	}

	if s.TraceID == "" {
		return nil, fmt.Errorf("trace ID is required")
	}
	now := time.Now()
	if l.location != nil {
		now = now.In(l.location)
	}
	l.observer.Dispatch(
		model.IngestionEvent{
			ID:        buildID(nil),
			Type:      model.IngestionEventTypeSpanUpdate,
			Timestamp: now,
			Body:      s,
		},
	)

	return s, nil
}

// Event 构建事件
func (l *Langfuse) Event(e *model.Event, parentID *string) (*model.Event, error) {
	if e.TraceID == "" {
		traceID, err := l.createTrace(e.Name)
		if err != nil {
			return nil, err
		}

		e.TraceID = traceID
	}

	e.ID = buildID(&e.ID)

	if parentID != nil {
		e.ParentObservationID = *parentID
	}
	now := time.Now()
	if l.location != nil {
		now = now.In(l.location)
	}
	l.observer.Dispatch(
		model.IngestionEvent{
			ID:        uuid.New().String(),
			Type:      model.IngestionEventTypeEventCreate,
			Timestamp: now,
			Body:      e,
		},
	)

	return e, nil
}

func (l *Langfuse) createTrace(traceName string) (string, error) {
	trace, errTrace := l.Trace(
		&model.Trace{
			Name: traceName,
		},
	)
	if errTrace != nil {
		return "", errTrace
	}

	return trace.ID, fmt.Errorf("unable to get trace ID")
}

// Flush 资源清理，等待所有观察者事件发送完成
func (l *Langfuse) Flush(ctx context.Context) {
	l.observer.Wait(ctx)
}

// buildID 构建ID
func buildID(id *string) string {
	if id == nil {
		return uuid.New().String()
	} else if *id == "" {
		return uuid.New().String()
	}

	return *id
}
