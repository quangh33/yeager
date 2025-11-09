package tracer

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
	"yeager/pkg/model"
)

type Tracer struct {
	CollectorURL string
	ServiceName  string
	Client       *http.Client

	queue chan *Span
	wg    sync.WaitGroup
	stop  chan struct{}

	// Atomic flag to track if the tracer is active.
	isRunning atomic.Bool
}

func NewTracer(collectorURL, serviceName string) *Tracer {
	t := &Tracer{
		CollectorURL: collectorURL,
		ServiceName:  serviceName,
		Client:       &http.Client{Timeout: 5 * time.Second},
		queue:        make(chan *Span, 1000),
		stop:         make(chan struct{}),
	}

	t.isRunning.Store(true)
	t.wg.Add(1)
	go t.backgroundReporter()
	return t
}

func (t *Tracer) Close() {
	t.isRunning.Store(false)
	close(t.stop)
	t.wg.Wait()
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

func (t *Tracer) backgroundReporter() {
	defer t.wg.Done()

	for {
		select {
		case span := <-t.queue:
			if err := t.send(span); err != nil {
				log.Printf("[Tracer] Failed to send span: %v", err)
			}
		case <-t.stop: // blocked until channel stop is closed
			for {
				select {
				case span := <-t.queue:
					if err := t.send(span); err != nil {
						log.Printf("[Tracer] Failed to send span during shutdown: %v", err)
					}
				default:
					return
				}
			}
		}
	}
}

func (t *Tracer) send(s_in *Span) error {
	span := model.Span{
		TraceID:       s_in.TraceID,
		SpanID:        s_in.SpanID,
		ParentSpanID:  s_in.ParentSpanID,
		OperationName: s_in.OperationName,
		StartTime:     s_in.StartTime,
		EndTime:       s_in.EndTime,
		Tags:          s_in.Tags,
		Process:       &model.Process{ServiceName: t.ServiceName},
	}

	data, err := json.Marshal(span)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/api/spans", t.CollectorURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := t.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("collector rejected span: status %d", resp.StatusCode)
	}

	return nil
}

func generateRandomID(bytesNum int) string {
	bytes := make([]byte, bytesNum)
	if _, err := rand.Read(bytes); err != nil {
		return fmt.Sprintf("%x", time.Now().UnixNano())
	}
	return hex.EncodeToString(bytes)
}
