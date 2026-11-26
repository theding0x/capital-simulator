-- +goose Up
-- Vol. I Ch. 15 seed fidelity (issue #450, v1-ch15 item). Each seeded magnitude is
-- Marx's own, derived from his figures below, or left in place and declared illustrative.
--
-- Calibration (sec. 2, steam-plough passage). At a 100% rate of surplus-value, 66 men in
-- one hour "produce ... a value of 30 shillings" while their wages, 15 shillings,
-- "represent only their labour for half an hour". 66 men x 60 min = 3,960 LabourMinutes
-- = 30s = 360d, so 1d = 11 LabourMinutes, 1s = 132, 1 pound = 2,640.
-- Check: 15s = 1,980 min = 66 men x 30 min.
--
-- m05 steam-plough. Marx: "does as much work in one hour at a cost of three-pence, as 66
-- men at a cost of 15 shillings". His 3d is the cost of a steam horse-power-hour (the
-- Morton footnote on steam power, sec. 1), taken here wholesale as value transferred (the
-- carpet example's "we take no account of buildings, coal, &c." simplification). 3d = 33
-- LabourMinutes per ploughing hour; 00002 carried 720000 / 500 days / 1 unit = 1,440 per
-- hour (43.6x his figure). Now productive_power 8 (ploughing hours in the 8-hour reference
-- day), lifespan_days 5475, machine_value 33 x 8 x 5475 = 1445400 (pound 547 10s), so wear
-- is 264 a day and 33 an hour. The plough's price and life are assumed: Marx gives neither.
-- The 5,475 days (15 years x 365) are borrowed from his sec. 3B illustration of "a machine"
-- worked 8 hours daily for 15 years, not a plough's life, and since wear is fixed at 33 min
-- per ploughing hour the life only scales machine_value. hand_labour_per_unit stays
-- 66 x 60 = 3960.
--
-- m06 / f02. The 366 units / 4400 min were Mr. Baynes' sec. 2 spinning figures (450 mule
-- spindles), seeded as a power-loom. It is now that mule in Marx's units: one unit is one
-- spindle's ten-hour yield of 13 oz, which a hand-spinner needs "sixty hours" = 3600 min
-- for, so productive_power 450 and 450 x 3600 = 1620000 min = his 27,000 hours. Its mill
-- takes his one horse-power and 2.5 operatives for a ten-hour day = 5 x 300 = 1500 min, half
-- necessary (the project's 100% convention; Marx's 100% is a hypothetical in the plough
-- passage). His text is ambiguous: each spindle, "working ten hours, produces 13 ounces",
-- yet the 2.5 operatives "spin weekly 365 5/8 lbs", and his 150 h is 2.5 operatives x six
-- ten-hour days, a weekly figure, while his 27,000 h is for the same 365 5/8 lb. The seed
-- takes the ten-hour reading, so one tick is one day (25 operative-hours against 27,000
-- hand-hours). His two figures, 27,000 h and 150 h, give 180x, which holds only for his
-- 60-hour week; the per-day factor is 1620000 / (240 + 1500) = 931, a modelling consequence
-- of that daily reading, not his figure.
--
-- Ids ...150m0N are not hex; machines become ...150a0N, with every factory_machines
-- reference (including rows added after seeding). Factory ids are unchanged.
--
-- Left as seeded and illustrative (Marx gives no price, life or horse-power): the
-- needle-machines' machine_value 600000 / 1000 days, the mule's 480000 / 2000 days, the
-- Needle Works' 30 hp and half-necessary day, the mill's prime_mover_kind 'steam' and its
-- name "Baynes Mule Mill" (Marx names Baynes only as a source and says only "one-horse
-- power"). His own: 145,000 needles in an 11-hour day, four machines to one attendant,
-- 660 minutes (11 h).

UPDATE machines SET id = '5eed00000000000000150a01' WHERE id = '5eed00000000000000150m01';
UPDATE machines SET id = '5eed00000000000000150a02' WHERE id = '5eed00000000000000150m02';
UPDATE machines SET id = '5eed00000000000000150a03' WHERE id = '5eed00000000000000150m03';
UPDATE machines SET id = '5eed00000000000000150a04' WHERE id = '5eed00000000000000150m04';
UPDATE machines SET id = '5eed00000000000000150a05', machine_value = 1445400, lifespan_days = 5475, productive_power = 8 WHERE id = '5eed00000000000000150m05';
UPDATE machines SET id = '5eed00000000000000150a06', name = 'Self-acting Mule (450 spindles)', working_tool = 'spindle', productive_power = 450, hand_labour_per_unit = 3600 WHERE id = '5eed00000000000000150m06';

UPDATE factory_machines SET machine_id = '5eed00000000000000150a01' WHERE machine_id = '5eed00000000000000150m01';
UPDATE factory_machines SET machine_id = '5eed00000000000000150a02' WHERE machine_id = '5eed00000000000000150m02';
UPDATE factory_machines SET machine_id = '5eed00000000000000150a03' WHERE machine_id = '5eed00000000000000150m03';
UPDATE factory_machines SET machine_id = '5eed00000000000000150a04' WHERE machine_id = '5eed00000000000000150m04';
UPDATE factory_machines SET machine_id = '5eed00000000000000150a05' WHERE machine_id = '5eed00000000000000150m05';
UPDATE factory_machines SET machine_id = '5eed00000000000000150a06' WHERE machine_id = '5eed00000000000000150m06';

UPDATE factories SET name = 'Baynes Mule Mill', prime_mover_horsepower = 1, attendants = 5, living_labour_per_attendant = 300, necessary_labour_per_attendant = 150 WHERE id = '5eed00000000000000150f02';

-- +goose Down
-- Restores 00002 / 00070: old ids, plough, power-loom row and the weaving-shed factory.
UPDATE factories SET name = 'Manchester Weaving Shed', prime_mover_horsepower = 60, attendants = 1, living_labour_per_attendant = 480, necessary_labour_per_attendant = 240 WHERE id = '5eed00000000000000150f02';

UPDATE factory_machines SET machine_id = '5eed00000000000000150m01' WHERE machine_id = '5eed00000000000000150a01';
UPDATE factory_machines SET machine_id = '5eed00000000000000150m02' WHERE machine_id = '5eed00000000000000150a02';
UPDATE factory_machines SET machine_id = '5eed00000000000000150m03' WHERE machine_id = '5eed00000000000000150a03';
UPDATE factory_machines SET machine_id = '5eed00000000000000150m04' WHERE machine_id = '5eed00000000000000150a04';
UPDATE factory_machines SET machine_id = '5eed00000000000000150m05' WHERE machine_id = '5eed00000000000000150a05';
UPDATE factory_machines SET machine_id = '5eed00000000000000150m06' WHERE machine_id = '5eed00000000000000150a06';

UPDATE machines SET id = '5eed00000000000000150m01' WHERE id = '5eed00000000000000150a01';
UPDATE machines SET id = '5eed00000000000000150m02' WHERE id = '5eed00000000000000150a02';
UPDATE machines SET id = '5eed00000000000000150m03' WHERE id = '5eed00000000000000150a03';
UPDATE machines SET id = '5eed00000000000000150m04' WHERE id = '5eed00000000000000150a04';
UPDATE machines SET id = '5eed00000000000000150m05', machine_value = 720000, lifespan_days = 500, productive_power = 1 WHERE id = '5eed00000000000000150a05';
UPDATE machines SET id = '5eed00000000000000150m06', name = 'Power-loom', working_tool = 'shuttle', productive_power = 366, hand_labour_per_unit = 4400 WHERE id = '5eed00000000000000150a06';
