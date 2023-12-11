CREATE TABLE hnh_data.notification (
    user_id int REFERENCES hnh_data.user_profile ON DELETE CASCADE,
    message TEXT NOT NULL CONSTRAINT message_is_not_empty CHECK (length(message) > 0),
    PRIMARY KEY (user_id, message)
);
