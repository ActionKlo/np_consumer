-- +goose Up
CREATE TABLE test (
    name varchar
);

-- +goose StatementBegin
-- +goose StatementEnd

-- +goose Down
DROP TABLE test;

-- +goose StatementBegin
-- +goose StatementEnd
