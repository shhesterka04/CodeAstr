-- +goose Up
ALTER TABLE users
ADD COLUMN family_status text,
ADD COLUMN type_of_activity text,
ADD COLUMN notification_time int,
ADD COLUMN last_action_time timestamp;

-- +goose Down
ALTER TABLE users
DROP COLUMN family_status,
DROP COLUMN type_of_activity,
DROP COLUMN notification_time,
DROP COLUMN last_action_time;
