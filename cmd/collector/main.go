package main

import (
	"log"
	"net/http"
	"yeager/pkg/api"
	"yeager/pkg/storage/inmem"
)

func main() {
	memStore := inmem.NewStore()
	apiHandler := api.NewAPIHandler(memStore)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/spans", apiHandler.SubmitSpanHandler)
	mux.HandleFunc("/api/traces/", apiHandler.GetTraceHandler)
	fs := http.FileServer(http.Dir("./cmd/collector/static"))
	mux.Handle("/", fs)

	port := ":8080"
	log.Printf("Yeager Collector starting on port %s...", port)
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatalf("Server successfully failed to start: %v", err)
	}
}
