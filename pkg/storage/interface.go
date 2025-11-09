package storage

import (
	"context"
	"yeager/pkg/model"
)

type SpanWriter interface {
	WriteSpan(ctx context.Context, span *model.Span) error
}

type SpanReader interface {
	// GetTrace retrieves all spans associated with a given TraceID.
	GetTrace(ctx context.Context, traceID model.TraceID) (*model.Trace, error)
}

type Storage interface {
	SpanReader
	SpanWriter
}
