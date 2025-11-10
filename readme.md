![img.png](img/img.png)

# 🔭 Yeager - A Distributed Tracing System
This project is a simplified, educational implementation of a distributed tracing system like [CNCF Jaeger](www.jaegertracing.io/)
or [Google Dapper](https://research.google/pubs/dapper-a-large-scale-distributed-systems-tracing-infrastructure/),
built from scratch in Go.

# ✨ Features
- **Concurrent, Scalable Backend Collector**: A Go service that ingests spans via an HTTP API.
It uses a worker pool (goroutines) and a buffered channel to process and save spans asynchronously without blocking incoming requests.
- **Pluggable Storage Backend**: Built on a generic storage.Storage interface with two implementations:
  - In-Memory
  - Sqlite
- **Async Go Client SDK**: A library to instruct applications to buffer spans and send them to the collector without blocking
the main application thread 
- **Web UI & Visualization**: A vanilla HTML/CSS/JS frontend showing the traces and spans
  ![](img/trace_search.png)
- **Dependency Graph**: A dynamic map show how services interact.
  ![](img/dep_graph.png)

# Architecture
![Architecture](img/architecture.png)

# 🚀 Getting Started
## Option 1: Docker Compose
1. Prerequisites: Docker and Docker Compose must be installed.
2. Build & Run:
```bash
docker-compose up --build
```
3. View the UI: Open your browser to http://localhost:8080.

## Option 2: Local Development
1. Prerequisites: Go (version 1.21 or later) must be installed.
2. Run the collector
```bash
go run cmd/collector/main.go
```
3. View the UI: Open your browser to http://localhost:8080.

# ⚙️ How to Generate Traces
While the collector is running (using either Docker or local-dev), open a second terminal and run the example application.
```bash
go run cmd/example-app/main.go
```