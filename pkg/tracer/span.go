package tracer

import (
	"log"
	"time"
	"yeager/pkg/model"
)

type Span struct {
	tracer        *Tracer
	TraceID       model.TraceID
	SpanID        model.SpanID
	ParentSpanID  model.SpanID
	OperationName string
	StartTime     time.Time
	EndTime       time.Time
	Tags          map[string]string
}

func (s *Span) SetTag(key, value string) {
	s.Tags[key] = value
}

func (s *Span) Finish() {
	if !s.tracer.isRunning.Load() {
		// Tracer is shut down. Drop this span
		return
	}
	s.EndTime = time.Now()
	select {
	case s.tracer.queue <- s:
	default:
		log.Println("[Tracer] Queue full, dropping span")
	}
}
