package model

import "time"

// TraceID is a unique identifier for a complete trace.
type TraceID string

// SpanID is a unique identifier for a single operation.
type SpanID string

// Span represents a single operation.
type Span struct {
	TraceID TraceID `json:"trace_id"`
	SpanID  SpanID  `json:"span_id"`
	// ParentSpanID is empty if this is a root span.
	ParentSpanID  SpanID            `json:"parent_span_id,omitempty"`
	OperationName string            `json:"operation_name"`
	StartTime     time.Time         `json:"start_time"`
	EndTime       time.Time         `json:"end_time"`
	Tags          map[string]string `json:"tags,omitempty"`
	Logs          []Log             `json:"logs,omitempty"`
	// Process represents the microservice that generated this span.
	Process *Process `json:"process,omitempty"`
}

// Log represents a timestamped event within a Span.
type Log struct {
	Timestamp time.Time         `json:"timestamp"`
	Fields    map[string]string `json:"fields"`
}

// Process describes the service emitting the spans.
type Process struct {
	ServiceName string            `json:"service_name"`
	Tags        map[string]string `json:"tags,omitempty"`
}

// Trace is an aggregate of all spans with the same TraceID.
type Trace struct {
	TraceID TraceID `json:"trace_id"`
	Spans   []*Span `json:"spans"`
}
