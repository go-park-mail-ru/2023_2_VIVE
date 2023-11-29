CREATE SCHEMA csat_data;
SET search_path to csat_data;

CREATE TABLE IF NOT EXISTS question (
    id serial PRIMARY KEY CONSTRAINT id_is_positive CHECK (id > 0),
    "name" text NOT NULL CONSTRAINT name_is_not_empty CHECK (length("name") > 0),
    "text" text NOT NULL CONSTRAINT text_is_not_empty CHECK (length("text") > 0),
    UNIQUE("name")
);

CREATE TABLE IF NOT EXISTS answer (
    id serial PRIMARY KEY CONSTRAINT id_is_positive CHECK (id > 0),
    stars smallint NOT NULL CONSTRAINT stars_is_not_negative CHECK (stars >= 0),
    message text DEFAULT NULL CONSTRAINT message_is_not_empty CHECK (length(message) > 0),
    question_id int REFERENCES question ON DELETE CASCADE,
    created_at timestamptz DEFAULT now()
);

CREATE TABLE IF NOT EXISTS user_info (
    user_id int PRIMARY KEY CONSTRAINT id_is_positive CHECK (user_id > 0),
    last_request_at timestamptz DEFAULT now()
)

