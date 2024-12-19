-- +goose Up
ALTER TABLE users
ALTER COLUMN tg_id TYPE BIGINT;

-- +goose Down
ALTER TABLE users
ALTER COLUMN tg_id TYPE INT;