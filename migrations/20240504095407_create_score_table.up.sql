CREATE TABLE game_score
(
    user_id INTEGER NOT NULL,
    game_id INTEGER NOT NULL,
    score   INTEGER NOT NULL,
    PRIMARY KEY (user_id, game_id)
);

CREATE INDEX idx_game_score_user_id ON game_score (user_id);