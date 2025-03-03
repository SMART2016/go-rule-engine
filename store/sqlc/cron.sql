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