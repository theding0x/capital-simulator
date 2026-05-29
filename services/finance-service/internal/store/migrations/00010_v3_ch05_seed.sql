-- +goose Up
-- Vol. III Ch. 5 seed — worked economies in constant capital, so the economy
-- panel comes up populated on a fresh MySQL volume. Each lifts the rate of
-- profit s/(c+v) by economising c while s and v stay put.
--
-- Seed ID convention: 5eed00000000000000<CC><XX> where CC=05 (chapter) and XX
-- is the row sequence. These 22-char IDs never collide with the 24-char hex
-- profit.NewEconomyAnalysisID() output.
--
--   Row 0501 waste_reduction: a spinning mill, 2000c + 1000v + 1000s. Reclaiming
--     40 lb of waste cotton at £5 (= £200) cheapens c to £1,800; p' rises from
--     33.33% (3333 bp) to 35.71% (3571 bp), a gain of 238 bp.
--   Row 0502 shared_conditions: a factory, 1200c + 600v + 600s, sharing a
--     communal boiler-house with eleven others saves £100 of its own boiler
--     cost; p' rises from 33.33% (3333 bp) to 35.29% (3529 bp), a gain of 196 bp.
--   Row 0503 raw_material_saving: Marx's §I two-enterprise illustration —
--     11,000c + 1,000v + 1,000s economises £2,000 of constant capital (the gap
--     between the two firms); p' rises from 8 1/3% (833 bp) to 10% (1000 bp),
--     a gain of 167 bp.

INSERT INTO economy_analyses
    (id, kind, constant_capital, variable_capital, surplus_value, saving,
     old_profit_rate_bp, new_profit_rate_bp, profit_rate_gain_bp, created_at)
VALUES
    ('5eed000000000000000501', 'waste_reduction', 2000, 1000, 1000, 200,
     3333, 3571, 238, '2026-05-01 00:00:00.000000'),
    ('5eed000000000000000502', 'shared_conditions', 1200, 600, 600, 100,
     3333, 3529, 196, '2026-05-01 00:00:00.000000'),
    ('5eed000000000000000503', 'raw_material_saving', 11000, 1000, 1000, 2000,
     833, 1000, 167, '2026-05-01 00:00:00.000000');

-- +goose Down
DELETE FROM economy_analyses WHERE id IN (
    '5eed000000000000000501',
    '5eed000000000000000502',
    '5eed000000000000000503'
);
