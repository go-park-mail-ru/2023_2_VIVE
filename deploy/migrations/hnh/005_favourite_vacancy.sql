CREATE TABLE hnh_data.favourite_vacancy (
    user_id int REFERENCES hnh_data.user_profile ON DELETE RESTRICT,
    vacancy_id int REFERENCES hnh_data.vacancy ON DELETE RESTRICT,
    created_at timestamptz DEFAULT now(),
    PRIMARY KEY (user_id, vacancy_id)
)

---- create above / drop below ----

DROP TABLE IF EXISTS hnh_data.favourite_vacancy RESTRICT;
