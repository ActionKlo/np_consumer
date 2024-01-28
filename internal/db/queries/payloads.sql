-- name: GetPayloads :many
SELECT * FROM payloads;

-- name: GetPayloadsWithSettings :many
SELECT message_id, tracking_number, event_type, event_time, data, s.url
FROM payloads JOIN settings s ON s.settings_id = payloads.settings_id;

-- name: SavePayload :one
INSERT INTO payloads (
    message_id, tracking_number, event_id, event_type, event_time, data, receiver_id
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING message_id;