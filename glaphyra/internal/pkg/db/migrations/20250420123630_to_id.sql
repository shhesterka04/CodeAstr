-- +goose Up
ALTER TABLE users
    ADD COLUMN type_id INT REFERENCES users_types(id),
    ADD COLUMN style_id INT REFERENCES styles(id),
    ADD COLUMN zodiac_sign_id INT REFERENCES zodiac_signs(id);

UPDATE users
SET
    type_id = users_types.id
    FROM users_types
WHERE users.type = users_types.type;

UPDATE users
SET
    style_id = styles.id
    FROM styles
WHERE users.style = styles.style;

UPDATE users
SET
    zodiac_sign_id = zodiac_signs.id
    FROM zodiac_signs
WHERE users.zodiac_sign = zodiac_signs.sign;

ALTER TABLE users
DROP COLUMN type,
DROP COLUMN style,
DROP COLUMN zodiac_sign;


-- +goose Down
ALTER TABLE users
    ADD COLUMN type VARCHAR(10),
ADD COLUMN style TEXT DEFAULT 'Дружелюбный стиль',
ADD COLUMN zodiac_sign TEXT;

UPDATE users
SET
    type = users_types.type
    FROM users_types
WHERE users.type_id = users_types.id;

UPDATE users
SET
    style = styles.style
    FROM styles
WHERE users.style_id = styles.id;

UPDATE users
SET
    zodiac_sign = zodiac_signs.sign
    FROM zodiac_signs
WHERE users.zodiac_sign_id = zodiac_signs.id;

ALTER TABLE users
DROP COLUMN type_id,
DROP COLUMN style_id,
DROP COLUMN zodiac_sign_id;