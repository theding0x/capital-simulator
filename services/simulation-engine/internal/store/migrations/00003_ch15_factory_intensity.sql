-- +goose Up
-- Ch. 15 §3C — Intensification of labour. When the working day is
-- shortened under machinery, the remaining hours fill more densely.
-- The factor is a multiplicative scalar (≥ 1.0) on the hand-labour
-- displaced per machine per tick. Zero means "unset", treated as 1.0
-- by the tick path.

ALTER TABLE factories
    ADD COLUMN intensity_factor DOUBLE NOT NULL DEFAULT 0;

-- +goose Down
ALTER TABLE factories
    DROP COLUMN intensity_factor;
