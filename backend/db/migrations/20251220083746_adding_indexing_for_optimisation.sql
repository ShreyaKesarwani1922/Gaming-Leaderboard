-- +goose Up
-- +goose StatementBegin
-- Add indexes to leaderboard table
CREATE INDEX IF NOT EXISTS idx_leaderboard_total_score
    ON gaming.leaderboard(total_score DESC);

CREATE INDEX IF NOT EXISTS idx_leaderboard_user_id
    ON gaming.leaderboard(user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Remove the indexes if we roll back
DROP INDEX IF EXISTS idx_leaderboard_total_score;
DROP INDEX IF EXISTS idx_leaderboard_user_id;
-- +goose StatementEnd
