-- +goose Up
-- Ch. 15 §2 — a factory's living labour. Machines transfer value (c) but only
-- the attendants' living labour creates new value (v + s). Per tick, each
-- attendant performs living_labour_per_attendant LabourMinutes, of which
-- necessary_labour_per_attendant reproduces their labour-power (v).

ALTER TABLE factories
    ADD COLUMN attendants                     BIGINT NOT NULL DEFAULT 0,
    ADD COLUMN living_labour_per_attendant    BIGINT NOT NULL DEFAULT 0,
    ADD COLUMN necessary_labour_per_attendant BIGINT NOT NULL DEFAULT 0;

-- Seed exemplars. Needle Works: Marx §8A, "one woman superintends four such
-- machines" in an 11-hour day (660 min, half necessary). Weaving Shed: one
-- weaver on the standard 8-hour reference day (480 min, half necessary).
UPDATE factories
    SET attendants = 1, living_labour_per_attendant = 660, necessary_labour_per_attendant = 330
    WHERE id = '5eed00000000000000150f01';
UPDATE factories
    SET attendants = 1, living_labour_per_attendant = 480, necessary_labour_per_attendant = 240
    WHERE id = '5eed00000000000000150f02';

-- +goose Down
ALTER TABLE factories
    DROP COLUMN necessary_labour_per_attendant,
    DROP COLUMN living_labour_per_attendant,
    DROP COLUMN attendants;
