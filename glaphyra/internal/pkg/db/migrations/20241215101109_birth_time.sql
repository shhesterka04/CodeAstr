-- +goose Up
ALTER TABLE users
ALTER COLUMN birth_time TYPE text;

-- +goose Down
ALTER TABLE users
ALTER COLUMN birth_time TYPE time;
