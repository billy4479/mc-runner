
-- name: GetMetadata :one
SELECT value FROM metadata
WHERE key = ?
LIMIT 1;
