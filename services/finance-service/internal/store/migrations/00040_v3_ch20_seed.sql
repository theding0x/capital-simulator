-- +goose Up
-- Vol. III Ch. 20 seeds — Marx's historical exemplars of merchant capital.
-- SubordinationIndex rises with stage (500/800 pre-capitalist → 3000 transition → 7000 subordinated),
-- honouring invariant 2 (inverse-development law).
-- Pre-capitalist rows (stage=1) have wage_labour=0 (invariant 1);
-- only the StageSubordinated row (stage=3) sets wage_labour=1, which is allowed.
INSERT INTO historical_merchant_capitals
    (id, name, stage, profit_source, wage_labour, subordination_index, created_at)
VALUES
    ('5eed000000000000002001', 'Venice carrying trade',             1, 'unequal_exchange',       0, 500,  '2026-05-01 00:00:00.000000'),
    ('5eed000000000000002002', 'Genoa merchant republic',           1, 'monopoly_trade',         0, 800,  '2026-05-01 00:00:01.000000'),
    ('5eed000000000000002003', 'Dutch carrying trade',              2, 'carrying_monopoly',      0, 3000, '2026-05-01 00:00:02.000000'),
    ('5eed000000000000002004', 'Subordination under industrial capital', 3, 'industrial_subordination', 1, 7000, '2026-05-01 00:00:03.000000');

-- +goose Down
DELETE FROM historical_merchant_capitals
WHERE id IN (
    '5eed000000000000002001',
    '5eed000000000000002002',
    '5eed000000000000002003',
    '5eed000000000000002004'
);
