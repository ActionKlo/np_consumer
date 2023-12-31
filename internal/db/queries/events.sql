-- name: CreateEvent :exec
INSERT INTO status_events (
    status_id, shipment_id, event_timestamp, event_description
) VALUES (
    $1, $2, $3, $4
)ON CONFLICT DO NOTHING;