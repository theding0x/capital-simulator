-- +goose Up
ALTER TABLE agents ADD COLUMN labour_minutes BIGINT NOT NULL DEFAULT 0;

-- +goose Down
ALTER TABLE agents DROP COLUMN labour_minutes;
