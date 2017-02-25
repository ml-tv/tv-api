
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

ALTER TABLE shows RENAME on_netflix TO is_on_netflix;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

