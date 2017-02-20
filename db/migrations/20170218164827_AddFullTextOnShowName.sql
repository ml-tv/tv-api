
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

ALTER TABLE shows ADD COLUMN name_vector tsvector;
CREATE INDEX tsv_name_idx ON shows USING gin(name_vector);
CREATE INDEX day_of_week_idx ON shows (day_of_week);
CREATE INDEX name_idx ON shows (name);

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION shows_vector_update() RETURNS trigger AS $shows_vector_update$
BEGIN
    IF TG_OP = 'INSERT' THEN
        NEW.name_vector = to_tsvector('pg_catalog.english', COALESCE(NEW.name, ''));
    END IF;
    IF TG_OP = 'UPDATE' THEN
        IF NEW.name <> OLD.name THEN
            NEW.name_vector = to_tsvector('pg_catalog.english', COALESCE(NEW.name, ''));
        END IF;
    END IF;
    RETURN NEW;
END
$shows_vector_update$ LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER shows_vector_update BEFORE INSERT OR UPDATE ON shows
    FOR EACH ROW EXECUTE PROCEDURE shows_vector_update();

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back