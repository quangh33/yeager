package sqlstore

import (
	"context"
	"database/sql"
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"yeager/pkg/model"
)

//go:embed schema.sql
var schemaSQL string

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) (*Store, error) {
	s := &Store{db: db}
	if err := s.initSchema(); err != nil {
		return nil, fmt.Errorf("failed to init schema: %w", err)
	}
	return s, nil
}

func (s *Store) initSchema() error {
	_, err := s.db.Exec(schemaSQL)
	return err
}

// WriteSpan inserts the span into the SQL database.
func (s *Store) WriteSpan(ctx context.Context, span *model.Span) error {
	data, err := json.Marshal(span)
	if err != nil {
		return fmt.Errorf("failed to marshal span: %w", err)
	}

	query := `
    INSERT INTO spans (trace_id, span_id, parent_id, operation_name, service_name, start_time, end_time, data)
    VALUES (?, ?, ?, ?, ?, ?, ?, ?)
    `

	_, err = s.db.ExecContext(ctx, query,
		span.TraceID,
		span.SpanID,
		span.ParentSpanID,
		span.OperationName,
		span.Process.ServiceName,
		span.StartTime,
		span.EndTime,
		string(data),
	)
	if err != nil {
		return fmt.Errorf("failed to insert span: %w", err)
	}

	return nil
}

func (s *Store) GetTrace(ctx context.Context, traceID model.TraceID) (*model.Trace, error) {
	query := `SELECT data FROM spans WHERE trace_id = ? ORDER BY start_time ASC`
	rows, err := s.db.QueryContext(ctx, query, traceID)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	trace := &model.Trace{
		TraceID: traceID,
		Spans:   make([]*model.Span, 0),
	}

	for rows.Next() {
		var data []byte
		if err := rows.Scan(&data); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		var span model.Span
		if err := json.Unmarshal(data, &span); err != nil {
			log.Printf("Error unmarshaling span data: %v", err)
			continue
		}
		trace.Spans = append(trace.Spans, &span)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(trace.Spans) == 0 {
		return nil, fmt.Errorf("trace not found")
	}

	return trace, nil
}

func (s *Store) GetDependencies(ctx context.Context) ([]model.DependencyLink, error) {
	query := `
    SELECT
        parent.service_name,
        child.service_name,
        COUNT(*)
    FROM spans AS child
    JOIN spans AS parent ON child.parent_id = parent.span_id
    WHERE child.service_name != parent.service_name
    GROUP BY parent.service_name, child.service_name
    `

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links []model.DependencyLink
	for rows.Next() {
		var link model.DependencyLink
		if err := rows.Scan(&link.Parent, &link.Child, &link.CallCount); err != nil {
			return nil, err
		}
		links = append(links, link)
	}
	return links, nil
}
