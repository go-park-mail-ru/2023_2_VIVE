CREATE TABLE IF NOT EXISTS hnh_data.vacancy_view (
    id serial PRIMARY KEY CONSTRAINT id_is_positive CHECK (id > 0),
    vacancy_id int REFERENCES hnh_data.vacancy ON DELETE CASCADE,
    applicant_id int REFERENCES hnh_data.applicant ON DELETE CASCADE,
    created_at timestamptz DEFAULT now(),
    UNIQUE (vacancy_id, applicant_id)
);

CREATE TABLE IF NOT EXISTS hnh_data.cv_view (
    id serial PRIMARY KEY CONSTRAINT id_is_positive CHECK (id > 0),
    cv_id int REFERENCES hnh_data.cv ON DELETE CASCADE,
    employer_id int REFERENCES hnh_data.employer ON DELETE CASCADE,
    created_at timestamptz DEFAULT now(),
    UNIQUE (cv_id, employer_id)
);

---- create above / drop below ----

DROP TABLE IF EXISTS hnh_data.vacancy_view CASCADE;

DROP TABLE IF EXISTS hnh_data.cv_view CASCADE;
