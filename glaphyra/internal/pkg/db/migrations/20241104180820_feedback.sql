-- +goose Up
CREATE TABLE feedback (
    tg_id INT,
    feedback TEXT,
    created_at TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS feedback;
