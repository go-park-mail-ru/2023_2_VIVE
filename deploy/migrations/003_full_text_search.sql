DROP TABLE IF EXISTS hnh_data.organization CASCADE ;

CREATE TEXT SEARCH DICTIONARY russian_ispell (
    TEMPLATE = ispell,
    DictFile = russian,
    AffFile = russian,
    StopWords = russian
);

CREATE TEXT SEARCH CONFIGURATION ru (COPY=russian);

ALTER TEXT SEARCH CONFIGURATION ru
ALTER MAPPING FOR hword, hword_part, word WITH russian_ispell, russian_stem;

SET default_text_search_config = 'ru';

ALTER TABLE hnh_data.vacancy 
    ADD COLUMN fts TSVECTOR,
    ADD COLUMN organization_name TEXT NOT NULL DEFAULT 'Название вашей компании'
        CONSTRAINT organization_name_is_not_empty CHECK (length(organization_name) > 0);

ALTER TABLE hnh_data.employer 
    ADD COLUMN organization_name TEXT NOT NULL DEFAULT 'Название вашей компании'
        CONSTRAINT organization_name_is_not_empty CHECK (length(organization_name) > 0),
    ADD COLUMN organization_description TEXT NOT NULL DEFAULT 'Описание компании'
        CONSTRAINT organization_description_is_not_empty CHECK (length(organization_description) > 0),
    DROP COLUMN organization_id;


CREATE INDEX vacancy_fts ON hnh_data.vacancy USING GIN (fts);

-- set weights for full text search in hnh_data.vacancy
CREATE OR REPLACE FUNCTION hnh_data.update_fts_column()
RETURNS TRIGGER AS $$
BEGIN
    IF (TG_OP = 'INSERT') OR
       (TG_OP = 'UPDATE' AND (
           NEW."name" IS DISTINCT FROM OLD."name" OR
           NEW.description IS DISTINCT FROM OLD.description OR
           NEW.organization_name IS DISTINCT FROM OLD.organization_name
       ))
    THEN
        NEW.fts = setweight(coalesce(to_tsvector(NEW."name"), ''), 'A') ||
                  setweight(coalesce(to_tsvector(NEW.description), ''), 'B') ||
                  setweight(coalesce(to_tsvector(NEW.organization_name), ''), 'C');
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE PLPGSQL;


CREATE OR REPLACE TRIGGER vacancy_fts_update
BEFORE INSERT OR UPDATE ON hnh_data.vacancy
FOR EACH ROW EXECUTE FUNCTION hnh_data.update_fts_column();
