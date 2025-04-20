-- +goose Up
CREATE TABLE users_types (
                             id SERIAL PRIMARY KEY,
                             type VARCHAR(30)
);

INSERT INTO users_types (type) VALUES ('user'), ('moderator');

CREATE TABLE languages (
                           id SERIAL PRIMARY KEY,
                           language VARCHAR(100)
);

INSERT INTO languages (language) VALUES ('Русский'), ('English');

CREATE TABLE styles (
                        id SERIAL PRIMARY KEY,
                        style VARCHAR(100)
);

INSERT INTO styles (style) VALUES ('Дружелюбный стиль'), ('Серьезный стиль'), ('Шутливый стиль');

CREATE TABLE zodiac_signs (
                              id SERIAL PRIMARY KEY,
                              sign VARCHAR(100)
);

INSERT INTO zodiac_signs (sign) VALUES ('Aris'), ('Taurus'), ('Gemini'), ('Cancer'), ('Leo'), ('Virgo'), ('Libra'), ('Scorpio'), ('Sagittarius'), ('Capricorn'), ('Aquarius'), ('Pisces');

ALTER TABLE users
    ADD COLUMN language_id INT REFERENCES languages(id);

-- +goose Down
DROP TABLE IF EXISTS users_types;
DROP TABLE IF EXISTS languages;
DROP TABLE IF EXISTS styles;
DROP TABLE IF EXISTS zodiac_signs;
ALTER TABLE users
    DROP COLUMN lang_id;
