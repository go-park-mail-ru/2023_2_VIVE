CREATE SCHEMA hnh_data AUTHORIZATION vive_admin;

SET search_path TO hnh_data;

DROP TABLE IF EXISTS user_profile CASCADE;

CREATE TABLE user_profile (
    id serial PRIMARY KEY
        CONSTRAINT id_is_positive CHECK (id > 0),
    email CHARACTER(256) NOT NULL UNIQUE,
    pswd bytea NOT NULL,
    salt bytea NOT NULL,
    first_name TEXT NOT NULL
        CONSTRAINT first_name_is_not_empty CHECK (length(first_name) > 0),
    last_name TEXT NOT NULL
        CONSTRAINT last_name_is_not_empty CHECK (length(last_name) > 0),
    birthday date DEFAULT NULL,
    phone_number CHARACTER(16) DEFAULT NULL
        CONSTRAINT phone_number_is_not_empty CHECK (length(phone_number) > 0),
    "location" TEXT DEFAULT NULL
        CONSTRAINT location_is_not_empty CHECK (length("location") > 0),
    avatar_path TEXT DEFAULT NULL
        CONSTRAINT avatar_path_is_not_empty CHECK (length(avatar_path) > 0)
);

DROP TABLE IF EXISTS organization CASCADE;

CREATE TABLE organization (
    id serial PRIMARY KEY
        CONSTRAINT id_is_positive CHECK (id > 0),
    "name" TEXT UNIQUE NOT NULL
        CONSTRAINT name_is_not_empty CHECK (length("name") > 0),
    location TEXT DEFAULT NULL
        CONSTRAINT location_is_not_empty CHECK (length("location") > 0),
    description TEXT NOT NULL
        CONSTRAINT description_is_not_empty CHECK (length(description) > 0)
);

DROP TABLE IF EXISTS employer CASCADE;

CREATE TABLE employer (
    id serial PRIMARY KEY
        CONSTRAINT id_is_positive CHECK (id > 0),
    organization_id int REFERENCES organization ON DELETE CASCADE,
    user_id int REFERENCES user_profile ON DELETE CASCADE,
    UNIQUE (user_id),
    UNIQUE (organization_id, user_id)
);

DROP TABLE IF EXISTS applicant CASCADE;

CREATE TABLE applicant (
    id serial PRIMARY KEY
        CONSTRAINT id_is_positive CHECK (id > 0),
    user_id int REFERENCES user_profile ON DELETE CASCADE,
    status TEXT NOT NULL DEFAULT 'searching'
        CONSTRAINT status_is_not_empty CHECK (length(status) > 0),
    UNIQUE (user_id)
);

DROP TABLE IF EXISTS vacancy CASCADE;

CREATE TABLE vacancy (
    id serial PRIMARY KEY
        CONSTRAINT id_is_positive CHECK (id > 0),
    employer_id int REFERENCES employer ON DELETE CASCADE,
    "name" TEXT NOT NULL
        CONSTRAINT name_is_not_empty CHECK (length("name") > 0),
    description TEXT NOT NULL
        CONSTRAINT description_is_not_empty CHECK (length(description) > 0),
    salary_lower_bound int DEFAULT NULL
        CONSTRAINT salary_lower_bound_is_positive CHECK (salary_lower_bound > 0),
    salary_upper_bound int DEFAULT NULL
        CONSTRAINT salary_upper_bound_is_positive CHECK (salary_upper_bound > 0),
    employment TEXT DEFAULT NULL
        CONSTRAINT employment_is_not_empty CHECK (length(employment) > 0),
    experience_lower_bound int DEFAULT NULL
        CONSTRAINT experience_lower_bound_is_not_negative CHECK (experience_lower_bound >= 0),
    experience_upper_bound int DEFAULT NULL
        CONSTRAINT experience_upper_bound_is_not_negative CHECK (experience_upper_bound >= 0),
    education_type TEXT NOT NULL DEFAULT 'secondary'
        CONSTRAINT education_type_is_not_empty CHECK (length(education_type) > 0),
    "location" TEXT DEFAULT NULL
        CONSTRAINT location_is_not_empty CHECK (length("location") > 0),
    created_at timestamptz DEFAULT now(),
    updated_at timestamptz DEFAULT now(),
    CONSTRAINT valid_salaries CHECK (salary_lower_bound <= salary_upper_bound),
    CONSTRAINT valid_experience CHECK (experience_lower_bound <= salary_upper_bound)
);

DROP TABLE IF EXISTS cv CASCADE;

CREATE TABLE cv (
    id serial PRIMARY KEY
        CONSTRAINT id_is_positive CHECK (id > 0),
    applicant_id int REFERENCES applicant ON DELETE CASCADE,
    profession TEXT NOT NULL
        CONSTRAINT profession_is_not_empty CHECK (length(profession) > 0),
    description TEXT NOT NULL
        CONSTRAINT description_is_not_empty CHECK (length(description) > 0),
    status TEXT NOT NULL DEFAULT 'searching'
        CONSTRAINT status_is_not_empty CHECK (length(status) > 0),
    created_at timestamptz DEFAULT now(),
    updated_at timestamptz DEFAULT now()
);

DROP TABLE IF EXISTS response CASCADE;

CREATE TABLE response (
    id serial PRIMARY KEY
        CONSTRAINT id_is_positive CHECK (id > 0),
    vacancy_id int REFERENCES vacancy ON DELETE CASCADE,
    cv_id int REFERENCES cv ON DELETE CASCADE,
    created_at timestamptz DEFAULT now(),
    updated_at timestamptz DEFAULT now(),
    UNIQUE (vacancy_id, cv_id)
);

DROP TABLE IF EXISTS experience CASCADE;

CREATE TABLE experience (
    id serial PRIMARY KEY
        CONSTRAINT id_is_positive CHECK (id > 0),
    cv_id int REFERENCES cv ON DELETE CASCADE,
    organization_name TEXT NOT NULL
        CONSTRAINT organization_name_is_not_empty CHECK (length(organization_name) > 0),
    "position" TEXT NOT NULL
        CONSTRAINT position_is_not_empty CHECK (length("position") > 0),
    description TEXT NOT NULL
        CONSTRAINT description_is_not_empty CHECK (length(description) > 0),
    start_date date NOT NULL,
    end_date date DEFAULT NULL
);

DROP TABLE IF EXISTS "language" CASCADE;

CREATE TABLE "language" (
    id serial PRIMARY KEY
        CONSTRAINT id_is_positive CHECK (id > 0),
    "name" TEXT NOT NULL
        CONSTRAINT name_is_not_empty CHECK (length("name") > 0),
    "level" TEXT NOT NULL
        CONSTRAINT level_is_not_empty CHECK (length("level") > 0),
    UNIQUE ("name", "level")
);

DROP TABLE IF EXISTS cv_language_assign CASCADE;

CREATE TABLE cv_language_assign (
    cv_id int REFERENCES cv ON DELETE CASCADE,
    language_id int REFERENCES "language" ON DELETE CASCADE,
    PRIMARY KEY (cv_id, language_id)
);

DROP TABLE IF EXISTS institution CASCADE;

CREATE TABLE institution (
    id serial PRIMARY KEY
        CONSTRAINT id_is_positive CHECK (id > 0),
    "name" TEXT NOT NULL
        CONSTRAINT name_is_not_empty CHECK (length("name") > 0),
    education_level TEXT NOT NULL
        CONSTRAINT education_level_is_not_empty CHECK (length(education_level) > 0),
    UNIQUE ("name", education_level)
);

DROP TABLE IF EXISTS major_field CASCADE;

CREATE TABLE major_field (
    id serial PRIMARY KEY
        CONSTRAINT id_is_positive CHECK (id > 0),
    "name" TEXT NOT NULL UNIQUE
        CONSTRAINT name_is_not_empty CHECK (length("name") > 0)
);

DROP TABLE IF EXISTS institution_major_assign CASCADE;

CREATE TABLE institution_major_assign (
    id serial PRIMARY KEY
        CONSTRAINT id_is_positive CHECK (id > 0),
    institution_id int REFERENCES institution ON DELETE CASCADE,
    major_field_id int REFERENCES major_field ON DELETE CASCADE,
    UNIQUE (institution_id, major_field_id)
);

DROP TABLE IF EXISTS education CASCADE;

CREATE TABLE education (
    cv_id int REFERENCES cv ON DELETE CASCADE,
    institution_major_id int REFERENCES institution_major_assign ON DELETE CASCADE,
    graduation_year CHARACTER(4) NOT NULL
        CONSTRAINT valid_graduation_year CHECK (length(graduation_year) = 4),
    PRIMARY KEY (cv_id, institution_major_id)
);

DROP TABLE IF EXISTS skill CASCADE;

CREATE TABLE skill (
    id serial PRIMARY KEY
        CONSTRAINT id_is_positive CHECK (id > 0),
    "name" TEXT NOT NULL UNIQUE
        CONSTRAINT name_is_not_empty CHECK (length("name") > 0)
);

DROP TABLE IF EXISTS cv_skill_assign CASCADE;

CREATE TABLE cv_skill_assign (
    cv_id int REFERENCES cv ON DELETE CASCADE,
    skill_id int REFERENCES skill ON DELETE CASCADE,
    PRIMARY KEY (cv_id, skill_id)
);

DROP TABLE IF EXISTS vacancy_skill_assign CASCADE;

CREATE TABLE vacancy_skill_assign (
    vacancy_id int REFERENCES vacancy ON DELETE CASCADE,
    skill_id int REFERENCES skill ON DELETE CASCADE,
    PRIMARY KEY (vacancy_id, skill_id)
);
