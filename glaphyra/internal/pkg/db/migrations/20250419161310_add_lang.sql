-- +goose Up
ALTER TABLE users
ADD COLUMN lang_id INT REFERENCES languages(id);

-- +goose Down
ALTER TABLE users
DROP COLUMN lang_id;
