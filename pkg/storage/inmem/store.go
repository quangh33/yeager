package inmem

import (
	"context"
	"fmt"
	"sync"
	"yeager/pkg/model"
)

type Store struct {
	mu     sync.RWMutex
	traces map[model.TraceID]*model.Trace
}

func NewStore() *Store {
	return &Store{
		traces: make(map[model.TraceID]*model.Trace),
	}
}

func (s *Store) WriteSpan(ctx context.Context, span *model.Span) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	trace, exists := s.traces[span.TraceID]
	if !exists {
		trace = &model.Trace{
			TraceID: span.TraceID,
			Spans:   make([]*model.Span, 0),
		}
		s.traces[span.TraceID] = trace
	}

	trace.Spans = append(trace.Spans, span)
	return nil
}

func (s *Store) GetTrace(ctx context.Context, traceID model.TraceID) (*model.Trace, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	trace, exists := s.traces[traceID]
	if !exists {
		return nil, fmt.Errorf("trace not found: %s", traceID)
	}

	return trace, nil
}
