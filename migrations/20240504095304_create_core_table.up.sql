CREATE TABLE IF NOT EXISTS core
(
    user_id  INTEGER PRIMARY KEY,
    nickname VARCHAR(255) NOT NULL
);

CREATE INDEX idx_user_data_nickname ON core (nickname);