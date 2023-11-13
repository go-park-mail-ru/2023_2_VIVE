CREATE TABLE hnh_data.education_institution (
    id serial PRIMARY KEY CONSTRAINT id_is_positive CHECK (id > 0),
    cv_id int REFERENCES cv ON DELETE CASCADE,
    "name" TEXT NOT NULL CONSTRAINT name_is_not_empty CHECK (length("name") > 0),
    major_field TEXT NOT NULL CONSTRAINT major_field_is_not_empty CHECK (length(major_field) > 0),
    graduation_year CHARACTER(4) NOT NULL CONSTRAINT valid_graduation_year CHECK (length(graduation_year) = 4)
);

ALTER TABLE
    hnh_data.cv
ADD
    COLUMN education_level TEXT NOT NULL CONSTRAINT education_level_is_not_empty CHECK (length(education_level) > 0),
ADD
    COLUMN first_name TEXT NOT NULL CONSTRAINT first_name_is_not_empty CHECK (length(first_name) > 0),
ADD
    COLUMN last_name TEXT NOT NULL CONSTRAINT last_name_is_not_empty CHECK (length(last_name) > 0),
ADD
    COLUMN middle_name TEXT CONSTRAINT middle_name_is_not_empty CHECK (length(middle_name) > 0),
ADD
    COLUMN gender TEXT NOT NULL CONSTRAINT gender_is_not_empty CHECK (length(gender) > 0),
ADD
    COLUMN birthday date DEFAULT NULL,
ADD
    COLUMN "location" TEXT DEFAULT NULL CONSTRAINT location_is_not_empty CHECK (length("location") > 0),
ALTER COLUMN
    description DROP NOT NULL;


DROP TABLE hnh_data.education CASCADE;

DROP TABLE hnh_data.institution_major_assign CASCADE;

DROP TABLE hnh_data.education_institution CASCADE;

DROP TABLE hnh_data.major_field CASCADE;
