
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE INDEX created_at_idx ON shows (created_at);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

