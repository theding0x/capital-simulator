-- +goose Up
-- Vol. III Ch. 14 seed — one row per counteracting force kind (§I–§VI),
-- plus one scenario applying forces §I–§V to the §law rising-composition base.
-- Seed IDs: 5eed000000000000001401–1407.
-- These 22-char IDs never collide with the 24-char hex output of New*ID().
--
-- Base periods (from the Ch.13 §law seed): c=50/100/200/300/400, v=100, s'=100%,
-- s=100 each period; base p' = 6667/5000/3333/2500/2000 bp.
--
-- Uplift per period per force (δ = s·weight/(c+v), weights: §I/§II=500,
-- §III/§V=1000, §IV=300; §VI stock-capital contributes 0):
--   Period 1 (C=150): 333+333+666+200+666 = 2198 → 6667+2198 = 8865
--   Period 2 (C=200): 250+250+500+150+500 = 1650 → 5000+1650 = 6650
--   Period 3 (C=300): 166+166+333+100+333 = 1098 → 3333+1098 = 4431
--   Period 4 (C=400): 125+125+250+75+250  = 825  → 2500+825  = 3325
--   Period 5 (C=500): 100+100+200+60+200  = 660  → 2000+660  = 2660
-- Law persists: FinalProfitRate=2660 < InitialProfitRate=8865.

INSERT INTO counteracting_forces
    (id, kind, excluded_from_general_rate, note, created_at)
VALUES
    ('5eed000000000000001401', 'intensified_exploitation',  0, 'Marx §I: more surplus-value extracted from same labour-force',    '2026-05-01 00:00:00.000000'),
    ('5eed000000000000001402', 'depressed_wages',           0, 'Marx §II: wages forced below value of labour-power',              '2026-05-01 00:00:00.000000'),
    ('5eed000000000000001403', 'cheaper_constant_capital',  0, 'Marx §III: rising productivity cheapens elements of constant c',  '2026-05-01 00:00:00.000000'),
    ('5eed000000000000001404', 'relative_overpopulation',   0, 'Marx §IV: reserve army keeps wages low, new branches low-comp',   '2026-05-01 00:00:00.000000'),
    ('5eed000000000000001405', 'foreign_trade',             0, 'Marx §V: cheaper inputs and higher s'' via colonial trade',       '2026-05-01 00:00:00.000000'),
    ('5eed000000000000001406', 'stock_capital',             1, 'Marx §VI: stock/interest-bearing capital excluded from avg rate', '2026-05-01 00:00:00.000000');

INSERT INTO counteracting_scenarios
    (id, label, forces_json, base_trajectory_json, initial_profit_rate, final_profit_rate, modified_rates_json, created_at)
VALUES (
    '5eed000000000000001407',
    'All five industrial forces — law still holds',
    '[{"id":"5eed000000000000001401","kind":"intensified_exploitation","excluded_from_general_rate":false,"note":"Marx §I: more surplus-value extracted from same labour-force","created_at":"2026-05-01T00:00:00Z"},{"id":"5eed000000000000001402","kind":"depressed_wages","excluded_from_general_rate":false,"note":"Marx §II: wages forced below value of labour-power","created_at":"2026-05-01T00:00:00Z"},{"id":"5eed000000000000001403","kind":"cheaper_constant_capital","excluded_from_general_rate":false,"note":"Marx §III: rising productivity cheapens elements of constant c","created_at":"2026-05-01T00:00:00Z"},{"id":"5eed000000000000001404","kind":"relative_overpopulation","excluded_from_general_rate":false,"note":"Marx §IV: reserve army keeps wages low, new branches low-comp","created_at":"2026-05-01T00:00:00Z"},{"id":"5eed000000000000001405","kind":"foreign_trade","excluded_from_general_rate":false,"note":"Marx §V: cheaper inputs and higher s via colonial trade","created_at":"2026-05-01T00:00:00Z"}]',
    '{"id":"5eed000000000000001301","label":"The Law As Such — rising composition","surplus_value_rate":10000,"periods":[{"constant_capital":50,"variable_capital":100,"surplus_value":100,"profit_rate":6667},{"constant_capital":100,"variable_capital":100,"surplus_value":100,"profit_rate":5000},{"constant_capital":200,"variable_capital":100,"surplus_value":100,"profit_rate":3333},{"constant_capital":300,"variable_capital":100,"surplus_value":100,"profit_rate":2500},{"constant_capital":400,"variable_capital":100,"surplus_value":100,"profit_rate":2000}],"profit_rates":[6667,5000,3333,2500,2000],"created_at":"2026-05-01T00:00:00Z"}',
    8865,
    2660,
    '[8865,6650,4431,3325,2660]',
    '2026-05-01 00:00:00.000000'
);

-- +goose Down
DELETE FROM counteracting_scenarios WHERE id IN (
    '5eed000000000000001407'
);
DELETE FROM counteracting_forces WHERE id IN (
    '5eed000000000000001401',
    '5eed000000000000001402',
    '5eed000000000000001403',
    '5eed000000000000001404',
    '5eed000000000000001405',
    '5eed000000000000001406'
);
