-- DROP SCHEMA hnh_data;

CREATE SCHEMA hnh_data AUTHORIZATION vive_admin;

-- DROP SEQUENCE hnh_data.applicant_id_seq;

CREATE SEQUENCE hnh_data.applicant_id_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 2147483647
	START 1
	CACHE 1
	NO CYCLE;
-- DROP SEQUENCE hnh_data.cv_id_seq;

CREATE SEQUENCE hnh_data.cv_id_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 2147483647
	START 1
	CACHE 1
	NO CYCLE;
-- DROP SEQUENCE hnh_data.education_institution_id_seq;

CREATE SEQUENCE hnh_data.education_institution_id_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 2147483647
	START 1
	CACHE 1
	NO CYCLE;
-- DROP SEQUENCE hnh_data.employer_id_seq;

CREATE SEQUENCE hnh_data.employer_id_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 2147483647
	START 1
	CACHE 1
	NO CYCLE;
-- DROP SEQUENCE hnh_data.experience_id_seq;

CREATE SEQUENCE hnh_data.experience_id_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 2147483647
	START 1
	CACHE 1
	NO CYCLE;
-- DROP SEQUENCE hnh_data.institution_major_assign_id_seq;

CREATE SEQUENCE hnh_data.institution_major_assign_id_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 2147483647
	START 1
	CACHE 1
	NO CYCLE;
-- DROP SEQUENCE hnh_data.language_id_seq;

CREATE SEQUENCE hnh_data.language_id_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 2147483647
	START 1
	CACHE 1
	NO CYCLE;
-- DROP SEQUENCE hnh_data.major_field_id_seq;

CREATE SEQUENCE hnh_data.major_field_id_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 2147483647
	START 1
	CACHE 1
	NO CYCLE;
-- DROP SEQUENCE hnh_data.respond_id_seq;

CREATE SEQUENCE hnh_data.respond_id_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 2147483647
	START 1
	CACHE 1
	NO CYCLE;
-- DROP SEQUENCE hnh_data.skill_id_seq;

CREATE SEQUENCE hnh_data.skill_id_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 2147483647
	START 1
	CACHE 1
	NO CYCLE;
-- DROP SEQUENCE hnh_data.user_profile_id_seq;

CREATE SEQUENCE hnh_data.user_profile_id_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 2147483647
	START 1
	CACHE 1
	NO CYCLE;
-- DROP SEQUENCE hnh_data.vacancy_id_seq;

CREATE SEQUENCE hnh_data.vacancy_id_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 2147483647
	START 1
	CACHE 1
	NO CYCLE;-- hnh_data.education_institution definition

-- Drop table

-- DROP TABLE hnh_data.education_institution;

CREATE TABLE hnh_data.education_institution (
	id serial4 NOT NULL,
	"name" text NOT NULL,
	education_level text NOT NULL,
	CONSTRAINT education_institution_name_education_level_key UNIQUE (name, education_level),
	CONSTRAINT education_institution_pkey PRIMARY KEY (id),
	CONSTRAINT education_level_is_not_empty CHECK ((length(education_level) > 0)),
	CONSTRAINT id_is_positive CHECK ((id > 0)),
	CONSTRAINT name_is_not_empty CHECK ((length(name) > 0))
);


-- hnh_data."language" definition

-- Drop table

-- DROP TABLE hnh_data."language";

CREATE TABLE hnh_data."language" (
	id serial4 NOT NULL,
	"name" text NOT NULL,
	"level" text NOT NULL,
	CONSTRAINT id_is_positive CHECK ((id > 0)),
	CONSTRAINT language_name_level_key UNIQUE (name, level),
	CONSTRAINT language_pkey PRIMARY KEY (id),
	CONSTRAINT level_is_not_empty CHECK ((length(level) > 0)),
	CONSTRAINT name_is_not_empty CHECK ((length(name) > 0))
);


-- hnh_data.major_field definition

-- Drop table

-- DROP TABLE hnh_data.major_field;

CREATE TABLE hnh_data.major_field (
	id serial4 NOT NULL,
	"name" text NOT NULL,
	CONSTRAINT id_is_positive CHECK ((id > 0)),
	CONSTRAINT major_field_name_key UNIQUE (name),
	CONSTRAINT major_field_pkey PRIMARY KEY (id),
	CONSTRAINT name_is_not_empty CHECK ((length(name) > 0))
);


-- hnh_data.skill definition

-- Drop table

-- DROP TABLE hnh_data.skill;

CREATE TABLE hnh_data.skill (
	id serial4 NOT NULL,
	"name" text NOT NULL,
	CONSTRAINT id_is_positive CHECK ((id > 0)),
	CONSTRAINT name_is_not_empty CHECK ((length(name) > 0)),
	CONSTRAINT skill_name_key UNIQUE (name),
	CONSTRAINT skill_pkey PRIMARY KEY (id)
);


-- hnh_data.user_profile definition

-- Drop table

-- DROP TABLE hnh_data.user_profile;

CREATE TABLE hnh_data.user_profile (
	id serial4 NOT NULL,
	email bpchar(256) NOT NULL,
	pswd bytea NOT NULL,
	first_name text NOT NULL,
	last_name text NOT NULL,
	birthday date NULL,
	phone_number bpchar(16) NULL DEFAULT NULL::bpchar,
	"location" text NULL,
	CONSTRAINT first_name_is_not_empty CHECK ((length(first_name) > 0)),
	CONSTRAINT id_is_positive CHECK ((id > 0)),
	CONSTRAINT last_name_is_not_empty CHECK ((length(last_name) > 0)),
	CONSTRAINT location_is_not_empty CHECK ((length(location) > 0)),
	CONSTRAINT phone_number_is_not_empty CHECK ((length(phone_number) > 0)),
	CONSTRAINT user_profile_email_key UNIQUE (email),
	CONSTRAINT user_profile_pkey PRIMARY KEY (id)
);


-- hnh_data.applicant definition

-- Drop table

-- DROP TABLE hnh_data.applicant;

CREATE TABLE hnh_data.applicant (
	id serial4 NOT NULL,
	user_id int4 NULL,
	status text NOT NULL DEFAULT 'serching'::text,
	CONSTRAINT applicant_pkey PRIMARY KEY (id),
	CONSTRAINT id_is_positive CHECK ((id > 0)),
	CONSTRAINT status_is_not_empty CHECK ((length(status) > 0)),
	CONSTRAINT applicant_user_id_fkey FOREIGN KEY (user_id) REFERENCES hnh_data.user_profile(id)
);


-- hnh_data.cv definition

-- Drop table

-- DROP TABLE hnh_data.cv;

CREATE TABLE hnh_data.cv (
	id serial4 NOT NULL,
	applicant_id int4 NULL,
	status text NOT NULL DEFAULT 'serching'::text,
	created_at timestamp NULL DEFAULT now(),
	updated_at timestamp NULL DEFAULT now(),
	CONSTRAINT cv_pkey PRIMARY KEY (id),
	CONSTRAINT id_is_positive CHECK ((id > 0)),
	CONSTRAINT status_is_not_empty CHECK ((length(status) > 0)),
	CONSTRAINT cv_applicant_id_fkey FOREIGN KEY (applicant_id) REFERENCES hnh_data.applicant(id)
);


-- hnh_data.cv_language_assign definition

-- Drop table

-- DROP TABLE hnh_data.cv_language_assign;

CREATE TABLE hnh_data.cv_language_assign (
	cv_id int4 NULL,
	language_id int4 NULL,
	CONSTRAINT cv_language_assign_cv_id_language_id_key UNIQUE (cv_id, language_id),
	CONSTRAINT cv_language_assign_cv_id_fkey FOREIGN KEY (cv_id) REFERENCES hnh_data.cv(id),
	CONSTRAINT cv_language_assign_language_id_fkey FOREIGN KEY (language_id) REFERENCES hnh_data."language"(id)
);


-- hnh_data.cv_skill_assign definition

-- Drop table

-- DROP TABLE hnh_data.cv_skill_assign;

CREATE TABLE hnh_data.cv_skill_assign (
	cv_id int4 NULL,
	skill_id int4 NULL,
	CONSTRAINT cv_skill_assign_cv_id_skill_id_key UNIQUE (cv_id, skill_id),
	CONSTRAINT cv_skill_assign_cv_id_fkey FOREIGN KEY (cv_id) REFERENCES hnh_data.cv(id),
	CONSTRAINT cv_skill_assign_skill_id_fkey FOREIGN KEY (skill_id) REFERENCES hnh_data.skill(id)
);


-- hnh_data.employer definition

-- Drop table

-- DROP TABLE hnh_data.employer;

CREATE TABLE hnh_data.employer (
	id serial4 NOT NULL,
	user_id int4 NULL,
	CONSTRAINT employer_pkey PRIMARY KEY (id),
	CONSTRAINT id_is_positive CHECK ((id > 0)),
	CONSTRAINT employer_user_id_fkey FOREIGN KEY (user_id) REFERENCES hnh_data.user_profile(id)
);


-- hnh_data.experience definition

-- Drop table

-- DROP TABLE hnh_data.experience;

CREATE TABLE hnh_data.experience (
	id serial4 NOT NULL,
	cv_id int4 NULL,
	organization_name text NOT NULL,
	description text NOT NULL,
	start_date date NOT NULL,
	end_date date NULL,
	CONSTRAINT description_is_not_empty CHECK ((length(description) > 0)),
	CONSTRAINT experience_pkey PRIMARY KEY (id),
	CONSTRAINT id_is_positive CHECK ((id > 0)),
	CONSTRAINT organization_name_is_not_empty CHECK ((length(organization_name) > 0)),
	CONSTRAINT experience_cv_id_fkey FOREIGN KEY (cv_id) REFERENCES hnh_data.cv(id)
);


-- hnh_data.institution_major_assign definition

-- Drop table

-- DROP TABLE hnh_data.institution_major_assign;

CREATE TABLE hnh_data.institution_major_assign (
	id serial4 NOT NULL,
	institution_id int4 NULL,
	major_field_id int4 NULL,
	CONSTRAINT id_is_positive CHECK ((id > 0)),
	CONSTRAINT institution_major_assign_institution_id_major_field_id_key UNIQUE (institution_id, major_field_id),
	CONSTRAINT institution_major_assign_pkey PRIMARY KEY (id),
	CONSTRAINT institution_major_assign_institution_id_fkey FOREIGN KEY (institution_id) REFERENCES hnh_data.education_institution(id),
	CONSTRAINT institution_major_assign_major_field_id_fkey FOREIGN KEY (major_field_id) REFERENCES hnh_data.major_field(id)
);


-- hnh_data.organization definition

-- Drop table

-- DROP TABLE hnh_data.organization;

CREATE TABLE hnh_data.organization (
	employer_id int4 NULL,
	"name" text NOT NULL,
	"location" text NULL,
	description text NOT NULL,
	CONSTRAINT description_is_not_empty CHECK ((length(description) > 0)),
	CONSTRAINT location_is_not_empty CHECK ((length(location) > 0)),
	CONSTRAINT name_is_not_empty CHECK ((length(name) > 0)),
	CONSTRAINT organization_name_key UNIQUE (name),
	CONSTRAINT organization_employer_id_fkey FOREIGN KEY (employer_id) REFERENCES hnh_data.employer(id)
);


-- hnh_data.vacancy definition

-- Drop table

-- DROP TABLE hnh_data.vacancy;

CREATE TABLE hnh_data.vacancy (
	id serial4 NOT NULL,
	employer_id int4 NULL,
	"name" text NOT NULL,
	description text NOT NULL,
	salary_lower_bound int4 NULL,
	salary_upper_bound int4 NULL,
	employment text NULL,
	experience_lower_bound int4 NULL,
	experience_upper_bound int4 NULL,
	education_type text NOT NULL DEFAULT 'secondary'::text,
	"location" text NULL,
	created_at timestamp NULL DEFAULT now(),
	updated_at timestamp NULL DEFAULT now(),
	CONSTRAINT description_is_not_empty CHECK ((length(description) > 0)),
	CONSTRAINT education_type_is_not_empty CHECK ((length(education_type) > 0)),
	CONSTRAINT employment_is_not_empty CHECK ((length(employment) > 0)),
	CONSTRAINT experience_lower_bound_is_positive CHECK ((salary_lower_bound > 0)),
	CONSTRAINT experience_upper_bound_is_positive CHECK ((salary_upper_bound > 0)),
	CONSTRAINT id_is_positive CHECK ((id > 0)),
	CONSTRAINT location_is_not_empty CHECK ((length(location) > 0)),
	CONSTRAINT name_is_not_empty CHECK ((length(name) > 0)),
	CONSTRAINT salary_lower_bound_is_positive CHECK ((salary_lower_bound > 0)),
	CONSTRAINT salary_upper_bound_is_positive CHECK ((salary_upper_bound > 0)),
	CONSTRAINT vacancy_name_key UNIQUE (name),
	CONSTRAINT vacancy_pkey PRIMARY KEY (id),
	CONSTRAINT valid_experience CHECK ((experience_lower_bound <= salary_upper_bound)),
	CONSTRAINT valid_salaries CHECK ((salary_lower_bound <= salary_upper_bound)),
	CONSTRAINT vacancy_employer_id_fkey FOREIGN KEY (employer_id) REFERENCES hnh_data.employer(id)
);


-- hnh_data.vacancy_skill_assign definition

-- Drop table

-- DROP TABLE hnh_data.vacancy_skill_assign;

CREATE TABLE hnh_data.vacancy_skill_assign (
	vacancy_id int4 NULL,
	skill_id int4 NULL,
	CONSTRAINT vacancy_skill_assign_vacancy_id_skill_id_key UNIQUE (vacancy_id, skill_id),
	CONSTRAINT vacancy_skill_assign_skill_id_fkey FOREIGN KEY (skill_id) REFERENCES hnh_data.skill(id),
	CONSTRAINT vacancy_skill_assign_vacancy_id_fkey FOREIGN KEY (vacancy_id) REFERENCES hnh_data.vacancy(id)
);


-- hnh_data.education definition

-- Drop table

-- DROP TABLE hnh_data.education;

CREATE TABLE hnh_data.education (
	cv_id int4 NULL,
	institution_major_id int4 NULL,
	graduation_year bpchar(4) NOT NULL,
	CONSTRAINT valid_graduation_year CHECK ((length(graduation_year) = 4)),
	CONSTRAINT education_cv_id_fkey FOREIGN KEY (cv_id) REFERENCES hnh_data.cv(id),
	CONSTRAINT education_institution_major_id_fkey FOREIGN KEY (institution_major_id) REFERENCES hnh_data.institution_major_assign(id)
);


-- hnh_data.respond definition

-- Drop table

-- DROP TABLE hnh_data.respond;

CREATE TABLE hnh_data.respond (
	id serial4 NOT NULL,
	vacancy_id int4 NULL,
	cv_id int4 NULL,
	created_at timestamp NULL DEFAULT now(),
	updated_at timestamp NULL DEFAULT now(),
	CONSTRAINT id_is_positive CHECK ((id > 0)),
	CONSTRAINT respond_pkey PRIMARY KEY (id),
	CONSTRAINT respond_cv_id_fkey FOREIGN KEY (cv_id) REFERENCES hnh_data.cv(id),
	CONSTRAINT respond_vacancy_id_fkey FOREIGN KEY (vacancy_id) REFERENCES hnh_data.vacancy(id)
);
