-- +goose Up
ALTER TABLE users ADD COLUMN password_hash VARCHAR(256) NOT NULL DEFAULT '';

-- +goose Down
ALTER TABLE users DROP COLUMN password_hash;