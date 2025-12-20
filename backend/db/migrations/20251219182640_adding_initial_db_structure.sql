-- +goose Up
-- +goose StatementBegin

-- Ensure schema exists (safe if already created)
CREATE SCHEMA IF NOT EXISTS gaming;

-- Set schema for this migration
SET search_path TO gaming, public;

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    join_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE game_sessions (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    score INT NOT NULL,
    game_mode VARCHAR(50) NOT NULL,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_game_sessions_user
        FOREIGN KEY (user_id)
            REFERENCES users(id)
            ON DELETE CASCADE
);

CREATE TABLE leaderboard (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    total_score INT NOT NULL,
    rank INT,
    CONSTRAINT fk_leaderboard_user
        FOREIGN KEY (user_id)
            REFERENCES users(id)
            ON DELETE CASCADE
);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS gaming.leaderboard;
DROP TABLE IF EXISTS gaming.game_sessions;
DROP TABLE IF EXISTS gaming.users;

-- +goose StatementEnd
