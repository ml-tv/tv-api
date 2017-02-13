
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

ALTER TABLE users ADD COLUMN is_admin boolean NOT NULL DEFAULT FALSE;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

