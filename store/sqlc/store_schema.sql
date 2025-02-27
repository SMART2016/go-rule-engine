CREATE TABLE processed_events (
                                  id SERIAL PRIMARY KEY,
                                  tenant_id TEXT NOT NULL,
                                  event_type TEXT NOT NULL,
                                  occurred_at TIMESTAMP NOT NULL
);
