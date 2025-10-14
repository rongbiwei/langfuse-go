package model

import "time"

// IngestionEventType 类型
type IngestionEventType string

const (
	IngestionEventTypeTraceCreate      = "trace-create"
	IngestionEventTypeGenerationCreate = "generation-create"
	IngestionEventTypeGenerationUpdate = "generation-update"
	IngestionEventTypeScoreCreate      = "score-create"
	IngestionEventTypeSpanCreate       = "span-create"
	IngestionEventTypeSpanUpdate       = "span-update"
	IngestionEventTypeEventCreate      = "event-create"
)

// IngestionEvent 事件
type IngestionEvent struct {
	Type      IngestionEventType `json:"type"`
	ID        string             `json:"id"`
	Timestamp time.Time          `json:"timestamp"`
	Metadata  any
	Body      any `json:"body"`
	FailCount int
}

// Trace 跟踪
type Trace struct {
	ID        string     `json:"id,omitempty"`
	Timestamp *time.Time `json:"timestamp,omitempty"`
	Name      string     `json:"name,omitempty"`
	UserID    string     `json:"userId,omitempty"`
	Input     any        `json:"input,omitempty"`
	Output    any        `json:"output,omitempty"`
	SessionID string     `json:"sessionId,omitempty"`
	Release   string     `json:"release,omitempty"`
	Version   string     `json:"version,omitempty"`
	Metadata  any        `json:"metadata,omitempty"`
	Tags      []string   `json:"tags,omitempty"`
	Public    bool       `json:"public,omitempty"`
}

// ObservationLevel 观察级别
type ObservationLevel string

const (
	ObservationLevelDebug   ObservationLevel = "DEBUG"
	ObservationLevelDefault ObservationLevel = "DEFAULT"
	ObservationLevelWarning ObservationLevel = "WARNING"
	ObservationLevelError   ObservationLevel = "ERROR"
)

// Generation 生成
type Generation struct {
	TraceID             string           `json:"traceId,omitempty"`
	Name                string           `json:"name,omitempty"`
	StartTime           *time.Time       `json:"startTime,omitempty"`
	Metadata            any              `json:"metadata,omitempty"`
	Input               any              `json:"input,omitempty"`
	Output              any              `json:"output,omitempty"`
	Level               ObservationLevel `json:"level,omitempty"`
	StatusMessage       string           `json:"statusMessage,omitempty"`
	ParentObservationID string           `json:"parentObservationId,omitempty"`
	Version             string           `json:"version,omitempty"`
	ID                  string           `json:"id,omitempty"`
	EndTime             *time.Time       `json:"endTime,omitempty"`
	CompletionStartTime *time.Time       `json:"completionStartTime,omitempty"`
	Model               string           `json:"model,omitempty"`
	ModelParameters     any              `json:"modelParameters,omitempty"`
	Usage               Usage            `json:"usage,omitempty"`
	PromptName          string           `json:"promptName,omitempty"`
	PromptVersion       int              `json:"promptVersion,omitempty"`
}

// Usage .
type Usage struct {
	Input      int       `json:"input,omitempty"`
	Output     int       `json:"output,omitempty"`
	Total      int       `json:"total,omitempty"`
	Unit       UsageUnit `json:"unit,omitempty"`
	InputCost  float64   `json:"inputCost,omitempty"`
	OutputCost float64   `json:"outputCost,omitempty"`
	TotalCost  float64   `json:"totalCost,omitempty"`

	PromptTokens           int `json:"promptTokens,omitempty"`
	CompletionTokens       int `json:"completionTokens,omitempty"`
	TotalTokens            int `json:"totalTokens,omitempty"`
	PromptCacheHitTokens   int `json:"promptCacheHitTokens,omitempty"`
	PromptCacheMissTokens  int `json:"promptCacheMissTokens,omitempty"`
	PromptCacheWriteTokens int `json:"promptCacheWriteTokens,omitempty"`

	CompletionAcceptedPredictionTokens int `json:"completionAcceptedPredictionTokens,omitempty"`
	CompletionAudioTokens              int `json:"completionAudioTokens,omitempty"`
	CompletionReasoningTokens          int `json:"completionReasoningTokens,omitempty"`
	CompletionRejectedPredictionTokens int `json:"completionRejectedPredictionTokens,omitempty"`
	CompletionCachedTokens             int `json:"completionCachedTokens,omitempty"`

	PromptAcceptedPredictionTokens int `json:"promptAcceptedPredictionTokens,omitempty"`
	PromptAudioTokens              int `json:"promptAudioTokens,omitempty"`
	PromptReasoningTokens          int `json:"promptReasoningTokens,omitempty"`
	PromptRejectedPredictionTokens int `json:"promptRejectedPredictionTokens,omitempty"`
	PromptCachedTokens             int `json:"promptCachedTokens,omitempty"`
}

// UsageUnit .
type UsageUnit string

const (
	ModelUsageUnitCharacters   UsageUnit = "CHARACTERS"
	ModelUsageUnitTokens       UsageUnit = "TOKENS"
	ModelUsageUnitMilliseconds UsageUnit = "MILLISECONDS"
	ModelUsageUnitSeconds      UsageUnit = "SECONDS"
	ModelUsageUnitImages       UsageUnit = "IMAGES"
)

// Score 分数
type Score struct {
	ID            string  `json:"id,omitempty"`
	TraceID       string  `json:"traceId,omitempty"`
	Name          string  `json:"name,omitempty"`
	Value         float64 `json:"value,omitempty"`
	ObservationID string  `json:"observationId,omitempty"`
	Comment       string  `json:"comment,omitempty"`
}

// Span 跨度
type Span struct {
	TraceID             string           `json:"traceId,omitempty"`
	Name                string           `json:"name,omitempty"`
	StartTime           *time.Time       `json:"startTime,omitempty"`
	Metadata            any              `json:"metadata,omitempty"`
	Input               any              `json:"input,omitempty"`
	Output              any              `json:"output,omitempty"`
	Level               ObservationLevel `json:"level,omitempty"`
	StatusMessage       string           `json:"statusMessage,omitempty"`
	ParentObservationID string           `json:"parentObservationId,omitempty"`
	Version             string           `json:"version,omitempty"`
	ID                  string           `json:"id,omitempty"`
	EndTime             *time.Time       `json:"endTime,omitempty"`
}

// Event 事件
type Event struct {
	TraceID             string           `json:"traceId,omitempty"`
	Name                string           `json:"name,omitempty"`
	StartTime           *time.Time       `json:"startTime,omitempty"`
	Metadata            any              `json:"metadata,omitempty"`
	Input               any              `json:"input,omitempty"`
	Output              any              `json:"output,omitempty"`
	Level               ObservationLevel `json:"level,omitempty"`
	StatusMessage       string           `json:"statusMessage,omitempty"`
	ParentObservationID string           `json:"parentObservationId,omitempty"`
	Version             string           `json:"version,omitempty"`
	ID                  string           `json:"id,omitempty"`
}

type M map[string]interface{}
