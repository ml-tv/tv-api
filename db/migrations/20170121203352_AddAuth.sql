
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE users (
  id UUID NOT NULL,
  created_at timestamptz NOT NULL,
  updated_at timestamptz NOT NULL,
  deleted_at timestamptz,
  email VARCHAR(255) NOT NULL UNIQUE,
  name VARCHAR(255) NOT NULL,
  password VARCHAR(255) NOT NULL,
  PRIMARY KEY (id)
);

CREATE TABLE sessions (
  id UUID NOT NULL,
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
  created_at timestamptz NOT NULL,
  updated_at timestamptz,
  deleted_at timestamptz,
  PRIMARY KEY (id)
);

-- +goose Down

DROP TABLE sessions, users;

-- SQL section 'Down' is executed when this migration is rolled back

