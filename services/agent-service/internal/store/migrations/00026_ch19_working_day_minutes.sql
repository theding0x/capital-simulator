-- +goose Up
-- Rename working_day_hours to working_day_minutes in wage_forms.
-- Existing seeds inserted hours; this migration renames the column.
-- The Go domain now stores LabourMinutes throughout Ch. 19.

ALTER TABLE wage_forms
    RENAME COLUMN working_day_hours TO working_day_minutes;

-- +goose Down
ALTER TABLE wage_forms
    RENAME COLUMN working_day_minutes TO working_day_hours;
