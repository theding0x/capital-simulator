-- +goose Up
-- Issue #221: parameterise DepartmentalCapital with a turnover number so the
-- keystone exchange I(v+s) <-> II(c) can settle on annualised flows when
-- Department I and Department II turn over at different speeds. 10000 bp = one
-- turnover per year (annual); existing rows default to annual, so their raw
-- and annualised figures coincide exactly as before.
ALTER TABLE departmental_capitals
    ADD COLUMN turnover_number_bp BIGINT NOT NULL DEFAULT 10000;

ALTER TABLE extended_departmental_capitals
    ADD COLUMN turnover_number_bp BIGINT NOT NULL DEFAULT 10000;

-- +goose Down
ALTER TABLE extended_departmental_capitals DROP COLUMN turnover_number_bp;
ALTER TABLE departmental_capitals DROP COLUMN turnover_number_bp;
