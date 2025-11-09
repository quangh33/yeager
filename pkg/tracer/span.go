package tracer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
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
	s.EndTime = time.Now()
	err := s.report()
	if err != nil {
		fmt.Printf("Error reporting span: %v\n", err)
	}
}

func (s *Span) report() error {
	modelSpan := model.Span{
		TraceID:       s.TraceID,
		SpanID:        s.SpanID,
		ParentSpanID:  s.ParentSpanID,
		OperationName: s.OperationName,
		StartTime:     s.StartTime,
		EndTime:       s.EndTime,
		Tags:          s.Tags,
		Process: &model.Process{
			ServiceName: s.tracer.ServiceName,
		},
	}

	data, err := json.Marshal(modelSpan)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/api/spans", s.tracer.CollectorURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.tracer.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("collector rejected span: status %d", resp.StatusCode)
	}

	return nil
}
