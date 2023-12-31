-- name: CreateAddress :exec
INSERT INTO addresses (
    address_id, country, street, city, zip_code
) VALUES (
    $1, $2, $3, $4, $5
)ON CONFLICT DO NOTHING;