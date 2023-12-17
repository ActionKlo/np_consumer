-- name: GetMessages :many
SELECT * FROM messages
ORDER BY id;

-- name: GetMessageByID :one
SELECT * FROM messages WHERE id = $1;

-- name: InsertMessage :one
INSERT INTO messages (
    id, time, sender, tracknumber, country, city, street, postcode
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) ON CONFLICT (id) DO UPDATE
SET time = excluded.time
RETURNING *;

-- name: InsertStatus :one
INSERT INTO statuses (
    id, messageid, status, time
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;