CREATE TABLE hnh_data.vacancy_responce_notification (
    user_id int REFERENCES hnh_data.user_profile ON DELETE CASCADE,
    vacancy_id int REFERENCES hnh_data.vacancy ON DELETE CASCADE,
    cv_id int REFERENCES hnh_data.cv ON DELETE CASCADE,
    message TEXT NOT NULL CONSTRAINT message_is_not_empty CHECK (length(message) > 0),
    created_at timestamptz DEFAULT now(),
    PRIMARY KEY (user_id, vacancy_id, cv_id)
);

---- create above / drop below ----
DROP TABLE IF EXISTS hnh_data.vacancy_responce_notification CASCADE;

