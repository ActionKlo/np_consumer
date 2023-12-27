-- +goose Up
CREATE TABLE addresses (
    address_id uuid primary key not null unique default gen_random_uuid(),
    country varchar not null,
    street varchar not null,
    city varchar not null,
    zip_code varchar not null
);

CREATE TABLE customers (
    customer_id uuid primary key not null unique default gen_random_uuid(),
    customer_address_id uuid not null references addresses(address_id),
    name varchar not null,
    last_name varchar not null,
    email varchar not null unique,
    phone_number varchar not null unique
);

CREATE TABLE senders (
    sender_id uuid primary key not null unique default gen_random_uuid(),
    sender_address_id uuid not null references addresses(address_id),
    name varchar not null,
    email varchar not null unique,
    phone_number integer not null unique
);


CREATE TABLE shipments (
    shipment_id uuid primary key not null unique default gen_random_uuid(),
    sender_id uuid not null references senders(sender_id),
    customer_id uuid not null references customers(customer_id),
    size varchar not null,
    weight real not null,
    count integer not null
);

CREATE TABLE status_events (
    status_id uuid primary key not null unique default gen_random_uuid(),
    shipment_id uuid not null references shipments(shipment_id),
    event_timestamp timestamp default now(),
    event_description varchar not null
);

-- +goose StatementBegin
-- +goose StatementEnd

-- +goose Down
DROP TABLE status_events;

DROP TABLE shipments;

DROP TABLE senders;

DROP TABLE customers;

DROP TABLE addresses;

-- +goose StatementBegin
-- +goose StatementEnd
