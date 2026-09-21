CREATE TABLE events (
    event_id VARCHAR PRIMARY KEY,
    event_type VARCHAR(100) NOT NULL,
    status VARCHAR(50) NOT NULL,
    parent_id VARCHAR NOT NULL,
    parent_type VARCHAR(100) NOT NULL,
    parent_metadata JSONB NOT NULL,
    timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
