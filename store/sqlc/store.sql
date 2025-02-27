-- name: SaveEvent :exec
INSERT INTO processed_events (tenant_id, event_type, occurred_at)
VALUES ($1, $2, NOW());

-- name: IsDuplicate :one
SELECT COUNT(*) FROM processed_events
WHERE tenant_id = $1
  AND event_type = $2
  AND occurred_at >= NOW() - INTERVAL '1 hour' * $3;

-- name: CleanupOldEvents :exec
DELETE FROM processed_events WHERE occurred_at < NOW() - INTERVAL '1 hour' * $1;
