CREATE SCHEMA hnh_data AUTHORIZATION vive_admin;

CREATE TABLE user_profile (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	email CHARACTER(256) NOT NULL UNIQUE,
	pswd bytea NOT NULL,
	salt bytea NOT NULL UNIQUE,
	first_name TEXT NOT NULL
		CONSTRAINT first_name_is_not_empty CHECK (length(first_name) > 0),
	last_name TEXT NOT NULL
		CONSTRAINT last_name_is_not_empty CHECK (length(last_name) > 0),
	birthday date DEFAULT NULL,
	phone_number CHARACTER(16) DEFAULT NULL
		CONSTRAINT phone_number_is_not_empty CHECK (length(phone_number) > 0),
	LOCATION TEXT DEFAULT NULL
		CONSTRAINT location_is_not_empty CHECK (length(LOCATION) > 0)
);

CREATE TABLE employer (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	user_id uuid REFERENCES user_profile ON DELETE CASCADE 
);

CREATE TABLE applicant (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	user_id uuid REFERENCES user_profile ON DELETE CASCADE,
	status TEXT NOT NULL DEFAULT 'serching'
		CONSTRAINT status_is_not_empty CHECK (length(status) > 0)
);

CREATE TABLE organization (
	employer_id uuid REFERENCES employer ON DELETE CASCADE,
	name TEXT UNIQUE NOT NULL
		CONSTRAINT name_is_not_empty CHECK (length(name) > 0),
	LOCATION TEXT DEFAULT NULL
		CONSTRAINT location_is_not_empty CHECK (length(LOCATION) > 0),
	description TEXT NOT NULL
		CONSTRAINT description_is_not_empty CHECK (length(description) > 0)
);

CREATE TABLE vacancy (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	employer_id uuid REFERENCES employer ON DELETE CASCADE,
	name TEXT UNIQUE NOT NULL
		CONSTRAINT name_is_not_empty CHECK (length(name) > 0),
	description TEXT NOT NULL
		CONSTRAINT description_is_not_empty CHECK (length(description) > 0),
	salary_lower_bound int DEFAULT NULL
		CONSTRAINT salary_lower_bound_is_positive CHECK (salary_lower_bound > 0),
	salary_upper_bound int DEFAULT NULL
		CONSTRAINT salary_upper_bound_is_positive CHECK (salary_upper_bound > 0),
	employment TEXT DEFAULT NULL
		CONSTRAINT employment_is_not_empty CHECK (length(employment) > 0),
	experience_lower_bound int DEFAULT NULL
		CONSTRAINT experience_lower_bound_is_positive CHECK (experience_lower_bound > 0),
	experience_upper_bound int DEFAULT NULL
		CONSTRAINT experience_upper_bound_is_positive CHECK (experience_upper_bound > 0),
	education_type TEXT NOT NULL DEFAULT 'secondary'
		CONSTRAINT education_type_is_not_empty CHECK (length(education_type) > 0),
	LOCATION TEXT DEFAULT NULL
		CONSTRAINT location_is_not_empty CHECK (length(LOCATION) > 0),
	created_at timestamp DEFAULT now(),
	updated_at timestamp DEFAULT now(),
	CONSTRAINT valid_salaries CHECK (salary_lower_bound <= salary_upper_bound),
	CONSTRAINT valid_experience CHECK (experience_lower_bound <= salary_upper_bound)
);

CREATE TABLE cv (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	applicant_id uuid REFERENCES applicant ON DELETE CASCADE,
	status TEXT NOT NULL DEFAULT 'serching'
		CONSTRAINT status_is_not_empty CHECK (length(status) > 0),
	created_at timestamp DEFAULT now(),
	updated_at timestamp DEFAULT now()
);

CREATE TABLE responce (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	vacancy_id uuid REFERENCES vacancy ON DELETE CASCADE,
	cv_id uuid REFERENCES cv ON DELETE CASCADE,
	created_at timestamp DEFAULT now(),
	updated_at timestamp DEFAULT now()
);

CREATE TABLE experience (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	cv_id uuid REFERENCES cv ON DELETE CASCADE,
	organization_name TEXT NOT NULL
		CONSTRAINT organization_name_is_not_empty CHECK (length(organization_name) > 0),
	description TEXT NOT NULL
		CONSTRAINT description_is_not_empty CHECK (length(description) > 0),
	start_date date NOT NULL,
	end_date date DEFAULT NULL
);

CREATE TABLE "language" (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	name TEXT NOT NULL
		CONSTRAINT name_is_not_empty CHECK (length(name) > 0),
	LEVEL TEXT NOT NULL
		CONSTRAINT level_is_not_empty CHECK (length(LEVEL) > 0),
	UNIQUE (name, LEVEL)
);

CREATE TABLE cv_language_assign (
	cv_id uuid REFERENCES cv ON DELETE CASCADE,
	language_id uuid REFERENCES "language" ON DELETE CASCADE,
	UNIQUE (cv_id, language_id)
);

CREATE TABLE education_institution (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	name TEXT NOT NULL
		CONSTRAINT name_is_not_empty CHECK (length(name) > 0),
	education_level TEXT NOT NULL
		CONSTRAINT education_level_is_not_empty CHECK (length(education_level) > 0),
	UNIQUE (name, education_level)
);

CREATE TABLE major_field (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	name TEXT NOT NULL UNIQUE
		CONSTRAINT name_is_not_empty CHECK (length(name) > 0)
);

CREATE TABLE institution_major_assign (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	institution_id uuid REFERENCES education_institution ON DELETE CASCADE,
	major_field_id uuid REFERENCES major_field ON DELETE CASCADE,
	UNIQUE (institution_id, major_field_id)
);

CREATE TABLE education (
    cv_id uuid REFERENCES cv ON DELETE CASCADE,
    institution_major_id uuid REFERENCES institution_major_assign ON DELETE CASCADE,
    graduation_year CHARACTER(4) NOT NULL
        CONSTRAINT valid_graduation_year CHECK (length(graduation_year) = 4)
);

CREATE TABLE skill (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL UNIQUE
        CONSTRAINT name_is_not_empty CHECK (length(name) > 0)
);

CREATE TABLE cv_skill_assign (
    cv_id uuid REFERENCES cv ON DELETE CASCADE,
    skill_id uuid REFERENCES skill ON DELETE CASCADE,
    UNIQUE (cv_id, skill_id)
);

CREATE TABLE vacancy_skill_assign (
    vacancy_id uuid REFERENCES vacancy ON DELETE CASCADE,
    skill_id uuid REFERENCES skill ON DELETE CASCADE,
    UNIQUE (vacancy_id, skill_id)
);
