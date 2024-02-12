-- +goose Up
CREATE TABLE receivers (
    receiver_id uuid not null unique primary key default gen_random_uuid(),
    url varchar not null
);


CREATE TABLE payloads (
    message_id uuid not null unique primary key,
    receiver_id uuid not null references receivers(receiver_id),
    tracking_number varchar not null,
    event_id uuid not null unique,
    event_type varchar not null,
    event_time timestamp not null default now(),
    data jsonb not null
);


INSERT INTO receivers (
    receiver_id, url
) VALUES (
    '10899528-d8a6-49c4-ab1f-2f02b98811dc', 'http://100.104.232.63:8084/c222e8eb-40c6-4b9b-9a9c-1eebdce3ae8c'
);

-- +goose StatementBegin
-- +goose StatementEnd

-- +goose Down
DROP TABLE payloads;

DROP TABLE receivers;

-- +goose StatementBegin
-- +goose StatementEnd

