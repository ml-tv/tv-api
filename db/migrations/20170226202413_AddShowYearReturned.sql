
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

ALTER TABLE shows ADD COLUMN year_released smallint NOT NULL;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

