-- +goose Up
-- Vol. III Ch. 8 seed — Marx's four worked production spheres, so the
-- compositions panel comes up populated on a fresh MySQL volume.
--
-- Seed ID convention: 5eed00000000000000<CC><XX> where CC=08 (chapter) and XX is
-- the row sequence. These 22-char IDs never collide with the 24-char hex
-- avgprofit.NewProductionSphereID() output.
--
-- Derived columns (computed once, hard-coded):
--
--   Sphere A (600c + 100v), s' = 100% (10000 bp):
--     constant_capital = 700 - 100 = 600
--     surplus_value    = 10000 * 100 / 10000 = 100
--     individual_profit_rate = round_half_up(100*10000 / 700)
--                            = (1000000 + 350) / 700 = 1000350 / 700 = 1429
--     variable_percent = 100*100 / 700 = 14
--     labour_power_index = 100*3600*100 / 700 = 51428
--
--   Sphere B (100c + 600v), s' = 100% (10000 bp):
--     constant_capital = 700 - 600 = 100
--     surplus_value    = 10000 * 600 / 10000 = 600
--     individual_profit_rate = round_half_up(600*10000 / 700)
--                            = (6000000 + 350) / 700 = 6000350 / 700 = 8571
--     variable_percent = 600*100 / 700 = 85
--     labour_power_index = 600*3600*100 / 700 = 308571
--
--   European capital (84c + 16v), s' = 100% (10000 bp):
--     constant_capital = 100 - 16 = 84
--     surplus_value    = 10000 * 16 / 10000 = 16
--     individual_profit_rate = round_half_up(16*10000 / 100)
--                            = (160000 + 50) / 100 = 1600
--     variable_percent = 16*100 / 100 = 16
--     labour_power_index = 16*3600*100 / 100 = 57600
--
--   Asian capital (16c + 84v), s' = 25% (2500 bp):
--     constant_capital = 100 - 84 = 16
--     surplus_value    = 2500 * 84 / 10000 = 21
--     individual_profit_rate = round_half_up(21*10000 / 100)
--                            = (210000 + 50) / 100 = 2100
--     variable_percent = 84*100 / 100 = 84
--     labour_power_index = 84*3600*100 / 100 = 302400
--
-- Note: Asian capital shows a HIGHER individual p' (21%) than European (16%)
-- despite a LOWER s' (25% vs 100%) — because its far larger v/C sets so much
-- more living labour in motion per £100 of total capital.

INSERT INTO production_spheres
    (id, name, c, v, s_rate, constant_capital, surplus_value, individual_profit_rate, variable_percent, labour_power_index, created_at)
VALUES
    ('5eed000000000000000801', 'Sphere A (600c+100v)', 700, 100, 10000, 600, 100, 1429, 14, 51428,  '2026-05-01 00:00:00.000000'),
    ('5eed000000000000000802', 'Sphere B (100c+600v)', 700, 600, 10000, 100, 600, 8571, 85, 308571, '2026-05-01 00:00:00.000000'),
    ('5eed000000000000000803', 'European capital (84c+16v)', 100, 16, 10000, 84, 16, 1600, 16, 57600,  '2026-05-01 00:00:00.000000'),
    ('5eed000000000000000804', 'Asian capital (16c+84v)',    100, 84,  2500, 16, 21, 2100, 84, 302400, '2026-05-01 00:00:00.000000');

-- +goose Down
DELETE FROM production_spheres WHERE id IN (
    '5eed000000000000000801',
    '5eed000000000000000802',
    '5eed000000000000000803',
    '5eed000000000000000804'
);
