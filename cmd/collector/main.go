package main

import (
	"database/sql"
	"embed"
	"io/fs"
	"log"
	"net/http"
	collector "yeager/cmd/collector/processor"
	"yeager/pkg/api"
	"yeager/pkg/storage/sqlstore"

	_ "modernc.org/sqlite"
)

//go:embed static
var staticEmbed embed.FS

func main() {
	log.Println("🔌 Connecting to SQLite database...")
	db, err := sql.Open("sqlite",
		"file:traces.db?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	sqliteStore, err := sqlstore.NewStore(db)
	if err != nil {
		log.Fatalf("Failed to initialize SQL store: %v", err)
	}
	spanProcessor := collector.NewSpanProcessor(sqliteStore, 1000, 10)
	defer spanProcessor.Close()
	apiHandler := api.NewAPIHandler(sqliteStore, spanProcessor)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/spans", apiHandler.SubmitSpanHandler)
	mux.HandleFunc("/api/traces/", apiHandler.GetTraceHandler)
	mux.HandleFunc("/api/dependencies", apiHandler.GetDependenciesHandler)

	staticFS, _ := fs.Sub(staticEmbed, "static")
	mux.Handle("/", http.FileServer(http.FS(staticFS)))

	port := ":8080"
	log.Printf("Yeager Collector starting on port %s...", port)
	log.Printf("Simple Yeager UI available at http://localhost%s", port)
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatalf("Server successfully failed to start: %v", err)
	}
}
