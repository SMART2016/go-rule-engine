CREATE TABLE IF NOT EXISTS processed_events (
                                                id BIGSERIAL PRIMARY KEY,       -- Auto-incrementing primary key
                                                tenant_id VARCHAR(255) NOT NULL,
                                                event_type VARCHAR(255) NOT NULL,
                                                event_sha VARCHAR(255) NOT NULL,
                                                event_details json,
                                                occurred_at TIMESTAMP DEFAULT NOW(),
                                                actual_event_persistentce_time TIMESTAMP DEFAULT NOW()
);

-- Create indexes separately
CREATE UNIQUE INDEX idx_unique_processed_events ON processed_events (tenant_id, event_type, event_sha);
CREATE INDEX idx_processed_events_tenant_time ON processed_events (tenant_id, event_type, occurred_at DESC);
CREATE INDEX idx_processed_events_occurred_at ON processed_events (occurred_at);

-- Enabling auto vaccum parameters for events table
ALTER TABLE processed_events SET (
    autovacuum_vacuum_threshold = 5000,  -- Start vacuum when 5000+ rows are dead
    autovacuum_vacuum_scale_factor = 0.05,  -- Vacuum when 5% of table is dead tuples
    autovacuum_analyze_threshold = 1000,  -- Analyze after 1000 updates/deletes
    autovacuum_analyze_scale_factor = 0.02  -- Analyze when 2% of table changes
    );


CREATE EXTENSION IF NOT EXISTS pg_cron;
SELECT cron.schedule(
               'event_cleanup',  -- Job name
               '*/10 * * * *',   -- Runs every 10 minutes
               $$WITH rows_to_delete AS (
        SELECT tenant_id FROM processed_events
        WHERE occurred_at < NOW() - INTERVAL '1 hour' * 24  -- TTL = 24 hours
        LIMIT 5000
    )
    DELETE FROM processed_events USING rows_to_delete
    WHERE processed_events.ctid = rows_to_delete.ctid;$$
);

SELECT cron.schedule(
               'event_cleanup_vacuum',
               '*/30 * * * *',  -- Runs every 30 minutes
               $$VACUUM ANALYZE processed_events;$$
       );


