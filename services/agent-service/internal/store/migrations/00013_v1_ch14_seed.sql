-- +goose Up
-- Seed data for Capital Vol. I, Ch. 14 — Division of Labour and Manufacture.
--
-- Two manufactures are seeded to cover the chapter's principal exemplars:
--
--   1. Caslon Type-Foundry — serial, "splitting" origin. Embodies the §2
--      proportionality fixture exactly: "there are four founders and two
--      breakers to one rubber: the founder casts 2,000 type an hour, the
--      breaker breaks up 4,000, and the rubber polishes 8,000." Total
--      output rate balances at 8,000/hr per role-group.
--
--   2. Birmingham Carriage Works — heterogeneous, "combination" origin.
--      Marx's other two-fold-origin exemplar (§1): wheelwrights,
--      harness-makers, locksmiths, painters, fitters — diverse handicrafts
--      formerly independent, now annexed under one capitalist.
--
-- Both surface the §5 skill hierarchy: skilled roles (founder, wheelwright,
-- harness-maker, locksmith) appear alongside unskilled (breaker, rubber,
-- mould-runner, painter, fitter), each priced at its own labour-power
-- value.

-- 1. Two new capitalists who own the manufactures.
INSERT INTO agents (id, name, class, money_balance, hoarding, labour_minutes, created_at, updated_at) VALUES
    ('5eed00000000000000140100', 'Henry Caslon',     'capitalist', 500000, 0, 0, '2026-05-01 09:00:00.000000', '2026-05-01 09:00:00.000000'),
    ('5eed00000000000000140101', 'Matthew Boulton',  'capitalist', 750000, 0, 0, '2026-05-01 09:00:00.000000', '2026-05-01 09:00:00.000000');

-- 2. Manufactures.
INSERT INTO manufactures
    (id, name, capitalist_id, form, origin, individual_working_day_minutes, created_at) VALUES
    ('5eed00000000000000140f01',
     'Caslon Type-Foundry',
     '5eed00000000000000140100',
     'serial',
     'splitting',
     720,
     '2026-05-01 10:00:00.000000'),
    ('5eed00000000000000140f02',
     'Birmingham Carriage Works',
     '5eed00000000000000140101',
     'heterogeneous',
     'combination',
     720,
     '2026-05-01 10:30:00.000000');

-- 3. Detail roles.
--
-- Caslon Type-Foundry (§2 proportionality fixture): 4 founders @ 2,000/hr,
-- 2 breakers @ 4,000/hr, 1 rubber @ 8,000/hr. Rate × headcount = 8,000/hr
-- per role.
INSERT INTO detail_roles
    (manufacture_id, ordinal, role_name, skill_level, output_rate_per_hour, head_count, tool_name) VALUES
    ('5eed00000000000000140f01', 0, 'founder', 'skilled',   2000, 4, 'casting-mould'),
    ('5eed00000000000000140f01', 1, 'breaker', 'unskilled', 4000, 2, 'breaking-iron'),
    ('5eed00000000000000140f01', 2, 'rubber',  'unskilled', 8000, 1, 'polishing-stone');

-- Birmingham Carriage Works (§1 combination): five detail roles drawn
-- from formerly independent handicrafts, now subsumed under one
-- capitalist. Skilled trades retain their craft; finishing tasks
-- (mould-runner, painter) are the new unskilled positions (§5).
INSERT INTO detail_roles
    (manufacture_id, ordinal, role_name, skill_level, output_rate_per_hour, head_count, tool_name) VALUES
    ('5eed00000000000000140f02', 0, 'wheelwright',   'skilled',   4, 2, 'spokeshave'),
    ('5eed00000000000000140f02', 1, 'harness-maker', 'skilled',   4, 2, 'awl'),
    ('5eed00000000000000140f02', 2, 'locksmith',     'skilled',   4, 1, 'file'),
    ('5eed00000000000000140f02', 3, 'painter',       'unskilled', 4, 2, 'brush'),
    ('5eed00000000000000140f02', 4, 'fitter',        'unskilled', 4, 3, '');

-- +goose Down
DELETE FROM detail_roles WHERE manufacture_id IN (
    '5eed00000000000000140f01',
    '5eed00000000000000140f02'
);
DELETE FROM manufactures WHERE id IN (
    '5eed00000000000000140f01',
    '5eed00000000000000140f02'
);
DELETE FROM agents WHERE id IN (
    '5eed00000000000000140100',
    '5eed00000000000000140101'
);
