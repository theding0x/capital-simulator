-- +goose Up
-- Ch. 14 gives no formula for the productive power of the collective
-- labourer, so the factor is a recorded observation on the manufacture, not a
-- function of its role count or headcount (issue #450). 1 means no gain
-- asserted; existing rows keep it.
ALTER TABLE manufactures
    ADD COLUMN productive_power_factor DOUBLE NOT NULL DEFAULT 1;

-- +goose Down
ALTER TABLE manufactures DROP COLUMN productive_power_factor;
