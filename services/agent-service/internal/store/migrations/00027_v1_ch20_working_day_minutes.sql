-- +goose Up
-- Rename working_day_hours to working_day_minutes in time_wage_sessions.
-- Existing seeds inserted hours; this migration renames the column.
-- The Go domain now stores LabourMinutes throughout Ch. 20.

ALTER TABLE time_wage_sessions
    RENAME COLUMN working_day_hours TO working_day_minutes;

-- +goose Down
ALTER TABLE time_wage_sessions
    RENAME COLUMN working_day_minutes TO working_day_hours;
