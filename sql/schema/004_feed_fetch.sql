-- +goose up
ALTER TABLE feeds ADD last_fetched_at TIMESTAMP NOT NULL DEFAULT NOW();

-- +goose down
ALTER TABLE feeds DROP COLUMN last_fetched_at;

