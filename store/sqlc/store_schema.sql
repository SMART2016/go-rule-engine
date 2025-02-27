CREATE TABLE IF NOT EXISTS processed_events (
                                                id BIGSERIAL PRIMARY KEY,       -- Auto-incrementing primary key
                                                tenant_id VARCHAR(255) NOT NULL,
                                                event_type VARCHAR(255) NOT NULL,
                                                event_sha VARCHAR(255) NOT NULL,
                                                Created_at TIMESTAMP DEFAULT NOW()
);

-- Create indexes separately
CREATE INDEX idx_tenant_event ON processed_events (tenant_id, event_type, event_sha);
CREATE INDEX idx_occurred_at ON processed_events (Created_at);
