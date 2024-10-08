-- +goose Up
CREATE TABLE users (
    tg_id INT PRIMARY KEY,
    username TEXT,
    type VARCHAR(10) CHECK (type IN ('user', 'moderator')),
    style TEXT,
    gender TEXT,
    registration_date DATE,
    birth_date DATE,
    zodiac_sign TEXT,
    birth_time TIME,
    birth_place TEXT,
    friend_code TEXT,
    tokens INT
);

-- +goose Down
DROP TABLE IF EXISTS users;
