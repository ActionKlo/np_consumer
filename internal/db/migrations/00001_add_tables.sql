-- +goose Up
CREATE TABLE messages (
    ID  varchar unique primary key,
    Time timestamp not null,
    Sender varchar not null,
    TrackNumber varchar not null,
    Country varchar not null,
    City varchar not null,
    Street varchar not null,
    PostCode varchar not null
);

CREATE TABLE statuses (
    ID varchar unique primary key,
    MessageID varchar not null,
    Status varchar not null,
    Time timestamp not null,
    foreign key (MessageID) references messages(ID)
);

-- +goose StatementBegin
-- +goose StatementEnd

-- +goose Down
DROP TABLE statuses;

DROP TABLE messages;

-- +goose StatementBegin
-- +goose StatementEnd
