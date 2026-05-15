-- +goose Up
-- Add optional FK linking a piece-wage contract to the wage form it was derived
-- from. Nullable — existing contracts and direct-entry callers are unaffected.

ALTER TABLE piece_wages
    ADD COLUMN wage_form_id CHAR(24) NULL;

-- +goose Down
ALTER TABLE piece_wages DROP COLUMN wage_form_id;
