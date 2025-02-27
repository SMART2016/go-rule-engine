I-- name: SaveEvent :exec
INSERT INTO processed_events (tenant_id, event_type, event_sha, occurred_at)
VALUES ($1, $2, $3, NOW());

-- name: IsDuplicate :one
SELECT EXISTS (
    SELECT 1 FROM processed_events
    WHERE tenant_id = $1
      AND event_type = $2
      AND event_sha = $3
      AND occurred_at >= NOW() - INTERVAL '1 hour' * $4
);
-- name: CleanupOldEvents :exec
DELETE FROM processed_events
WHERE occurred_at < NOW() - INTERVAL '1 hour' * $1
    LIMIT 10000;  -- Prevents large deletions from locking the table
