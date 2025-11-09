package tracer

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"
	"yeager/pkg/model"
)

type Tracer struct {
	CollectorURL string
	ServiceName  string
	Client       *http.Client
}

func NewTracer(collectorURL, serviceName string) *Tracer {
	return &Tracer{
		CollectorURL: collectorURL,
		ServiceName:  serviceName,
		Client:       &http.Client{Timeout: 5 * time.Second},
	}
}

type spanContextKey struct{}

func (t *Tracer) StartSpan(ctx context.Context, operationName string) (context.Context, *Span) {
	var traceID model.TraceID
	var parentID model.SpanID

	if parent, ok := ctx.Value(spanContextKey{}).(*Span); ok {
		traceID = parent.TraceID
		parentID = parent.SpanID
	} else {
		traceID = model.TraceID(generateRandomID(16))
	}

	spanID := model.SpanID(generateRandomID(8))

	span := &Span{
		tracer:        t,
		TraceID:       traceID,
		SpanID:        spanID,
		ParentSpanID:  parentID,
		OperationName: operationName,
		StartTime:     time.Now(),
		Tags:          make(map[string]string),
	}

	newCtx := context.WithValue(ctx, spanContextKey{}, span)
	return newCtx, span
}

func generateRandomID(bytesNum int) string {
	bytes := make([]byte, bytesNum)
	if _, err := rand.Read(bytes); err != nil {
		return fmt.Sprintf("%x", time.Now().UnixNano())
	}
	return hex.EncodeToString(bytes)
}
