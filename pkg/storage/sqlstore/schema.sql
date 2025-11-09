CREATE TABLE IF NOT EXISTS spans (
     trace_id VARCHAR(64) NOT NULL,
     span_id VARCHAR(64) NOT NULL PRIMARY KEY,
     parent_id VARCHAR(64),
     operation_name VARCHAR(255) NOT NULL,
     start_time DATETIME NOT NULL,
     end_time DATETIME NOT NULL,
     data TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_trace_id ON spans(trace_id);
CREATE INDEX IF NOT EXISTS idx_operation_name ON spans(operation_name);
CREATE INDEX IF NOT EXISTS idx_start_time ON spans(start_time);