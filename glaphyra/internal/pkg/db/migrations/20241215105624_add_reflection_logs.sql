-- +goose Up
-- +goose StatementBegin
CREATE TABLE reflections (
     tg_id INT,
     created_at TIMESTAMP,
     mood_rating INT,
     emotions TEXT,
     activity TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS reflections;
-- +goose StatementEnd
