-- +goose Up
-- +goose StatementBegin

ALTER TABLE gaming.leaderboard
    ADD CONSTRAINT uq_leaderboard_user UNIQUE (user_id);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

ALTER TABLE gaming.leaderboard
DROP CONSTRAINT IF EXISTS uq_leaderboard_user;

-- +goose StatementEnd
