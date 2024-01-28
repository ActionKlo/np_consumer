-- +goose Up
CREATE TABLE settings (
    settings_id uuid not null unique primary key default gen_random_uuid(),
    receiver_id uuid not null unique,
    url varchar not null
);


CREATE TABLE payloads (
    message_id uuid not null unique primary key,
    tracking_number varchar not null,
    event_id uuid not null unique,
    event_type varchar not null,
    event_time timestamp not null default now(),
    data jsonb not null,
    receiver_id uuid not null
);


INSERT INTO settings (
    receiver_id, url
) VALUES (
    '10899528-d8a6-49c4-ab1f-2f02b98811dc', 'http://100.104.232.63:8084/8c5a5016-76ea-431e-8891-e1e60d4a274f'
);

-- +goose StatementBegin
-- +goose StatementEnd

-- +goose Down
DROP TABLE payloads;

DROP TABLE settings;

-- +goose StatementBegin
-- +goose StatementEnd

