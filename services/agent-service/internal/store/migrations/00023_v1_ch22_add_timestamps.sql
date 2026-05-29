-- +goose Up
-- Add created_at audit column to the two Ch.22 tables that were missing it.
-- All other domain tables have created_at DATETIME(6). Existing rows receive
-- CURRENT_TIMESTAMP as default.

ALTER TABLE national_intensities
    ADD COLUMN created_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6);

ALTER TABLE day_wages
    ADD COLUMN created_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6);

-- +goose Down
ALTER TABLE national_intensities DROP COLUMN created_at;
ALTER TABLE day_wages DROP COLUMN created_at;
