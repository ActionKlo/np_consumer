-- name: SaveSettings :one
INSERT INTO receivers (
    receiver_id, url
) VALUES (
    $1, $2
) RETURNING receiver_id;

-- name: GetSettingsByReceiverID :one
SELECT url FROM receivers where receiver_id = $1;

-- ----------------------------------

-- name: CreateReceiver :one
INSERT INTO receivers (
    receiver_id, url
) VALUES (
    $1, $2
) RETURNING receiver_id;

-- name: RetrieveReceiver :one
SELECT * FROM receivers WHERE receiver_id = $1;

-- name: DeleteReceiver :exec
DELETE FROM receivers WHERE receiver_id = $1;