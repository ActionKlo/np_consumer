-- name: CreateShipment :exec
INSERT INTO shipments (
    shipment_id, sender_id, customer_id, size, weight, count
) VALUES (
    $1, $2, $3, $4, $5, $6
) ON CONFLICT DO NOTHING;