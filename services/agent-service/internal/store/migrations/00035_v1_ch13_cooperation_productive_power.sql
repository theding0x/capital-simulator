-- +goose Up
-- Ch. 13 gives Marx no formula for the productive power of co-operation, so
-- the factor is a recorded observation on the cooperation, not a function of
-- its headcount (issue #450). 1 means no gain asserted; existing rows keep it.
ALTER TABLE cooperations
    ADD COLUMN productive_power_factor DOUBLE NOT NULL DEFAULT 1;

-- +goose Down
ALTER TABLE cooperations DROP COLUMN productive_power_factor;
