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
    ADD COLUMN fts TSVECTOR;

ALTER TABLE hnh_data.employer 
    ADD COLUMN organization_name TEXT NOT NULL DEFAULT 'Название вашей компании'
        CONSTRAINT organization_name_is_not_empty CHECK (length(organization_name) > 0),
    ADD COLUMN organization_description TEXT NOT NULL DEFAULT 'Описание компании'
        CONSTRAINT organization_description_is_not_empty CHECK (length(organization_description) > 0),
    DROP COLUMN organization_id;


CREATE INDEX vacancy_fts ON hnh_data.vacancy USING GIN (fts);

UPDATE
    hnh_data.vacancy
SET
    fts = setweight(to_tsvector("name"), 'A') 
        || setweight(to_tsvector(description), 'B')
        || setweight(to_tsvector(company_name), 'C'); 
