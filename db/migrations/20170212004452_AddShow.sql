
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE shows (
  id UUID NOT NULL,
  created_at timestamptz NOT NULL,
  updated_at timestamptz NOT NULL,
  deleted_at timestamptz,
  tmdb_id integer NOT NULL UNIQUE,
  name VARCHAR NOT NULL CHECK (length(name) < 100),
  original_name VARCHAR CHECK (length(original_name) < 100),
  synopsis VARCHAR NOT NULL CHECK (length(synopsis) < 10000),
  poster_path VARCHAR NOT NULL CHECK (length(poster_path) < 100),
  backdrop_path VARCHAR NOT NULL CHECK (length(backdrop_path) < 100),
  status int DEFAULT 0 NOT NULL CHECK (status < 4),
  day_of_week int DEFAULT 0 CHECK (day_of_week < 7),
  returning_date VARCHAR CHECK (length(returning_date) < 100),
  website VARCHAR CHECK (length(website) < 255),
  wikipedia VARCHAR CHECK (length(wikipedia) < 255),
  extra_link VARCHAR CHECK (length(extra_link) < 255),
  on_netflix BOOLEAN NOT NULL DEFAULT FALSE,
  PRIMARY KEY (id)
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

