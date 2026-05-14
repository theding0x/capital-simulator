-- +goose Up
-- Add optional FK linking a time-wage session to the wage form it was derived
-- from. Nullable — existing sessions and direct-entry callers are unaffected.

ALTER TABLE time_wage_sessions
    ADD COLUMN wage_form_id CHAR(24) NULL;

-- +goose Down
ALTER TABLE time_wage_sessions DROP COLUMN wage_form_id;
