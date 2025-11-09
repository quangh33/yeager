package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	collector "yeager/cmd/collector/processor"
	"yeager/pkg/model"
	"yeager/pkg/storage"
)

type Handler struct {
	storage   storage.Storage
	processor *collector.SpanProcessor
}

func NewAPIHandler(s storage.Storage, p *collector.SpanProcessor) *Handler {
	return &Handler{storage: s, processor: p}
}

// SubmitSpanHandler accepts a JSON representation of a Span and saves it.
func (h *Handler) SubmitSpanHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	var span model.Span
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&span); err != nil {
		http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
		return
	}

	if span.TraceID == "" || span.SpanID == "" {
		http.Error(w, "trace_id and span_id are required", http.StatusBadRequest)
		return
	}

	h.processor.ProcessSpan(&span)
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Span accepted"))
}

// GetTraceHandler retrieves a trace by its ID.
// Endpoint: GET /api/traces/{trace_id}
func (h *Handler) GetTraceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
		return
	}

	traceIDStr := strings.TrimPrefix(r.URL.Path, "/api/traces/")
	if traceIDStr == "" || traceIDStr == r.URL.Path {
		http.Error(w, "trace ID required", http.StatusBadRequest)
		return
	}

	traceID := model.TraceID(traceIDStr)

	trace, err := h.storage.GetTrace(r.Context(), traceID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Trace not found: %v", err), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(trace)
}

// GetDependenciesHandler returns the service dependency graph.
// Endpoint: GET /api/dependencies
func (h *Handler) GetDependenciesHandler(w http.ResponseWriter, r *http.Request) {
	links, err := h.storage.GetDependencies(r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get dependencies: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(links)
}
