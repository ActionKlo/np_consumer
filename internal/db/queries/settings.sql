-- name: SaveSettings :one
INSERT INTO settings (
    receiver_id, url
) VALUES (
    $1, $2
) RETURNING settings_id;

-- name: GetSettingsByReceiverID :one
SELECT url FROM settings where receiver_id = $1;