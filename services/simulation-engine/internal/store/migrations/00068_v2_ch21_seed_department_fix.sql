-- +goose Up
-- Fix: the Vol. II Ch. 21 extended-reproduction seed (00060) stored department
-- codes 'I' / 'II', but the CapitalDepartment enum — and every code path that
-- reads these rows — uses 'department_i' / 'department_ii'. The mismatch made
-- TickExtendedScheme's department lookup miss on every scheduler pass, so the
-- reproduction ticker errored each tick ("extended scheme 5eed…2101: machinery:
-- not found") and the seeded Year-I scheme never advanced. Memory-store unit
-- tests don't load SQL seeds, so CI stayed green — this only surfaced on a booted
-- MySQL stack. Align the seeded data with the enum (idempotent: it only touches
-- the legacy 'I'/'II' codes, so re-running over already-correct rows is a no-op).
UPDATE extended_departmental_capitals SET department = 'department_i'  WHERE department = 'I';
UPDATE extended_departmental_capitals SET department = 'department_ii' WHERE department = 'II';
UPDATE reinvestments                  SET department = 'department_i'  WHERE department = 'I';
UPDATE reinvestments                  SET department = 'department_ii' WHERE department = 'II';

-- +goose Down
-- One-way data correction aligning seed codes with the CapitalDepartment enum;
-- no meaningful rollback.
SELECT 1;
