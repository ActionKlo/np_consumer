-- name: CreateSender :exec
INSERT INTO senders (
    sender_id, sender_address_id, name, email, phone_number
) VALUES (
    $1, $2, $3, $4, $5
)ON CONFLICT DO NOTHING;