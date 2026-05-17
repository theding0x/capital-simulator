-- +goose Up
-- Ch. 30 seed: Westphalian flax spinning — the classic case of proletarianisation
-- creating both a free wage-labour force and a home market for manufactured goods.
-- Seed IDs follow 5eed00000000000000<CC><XX> where CC=30 (hex 1e).
INSERT INTO domestic_industries (id, historical_stage_id, name, households_engaged, annual_output_pence, created_at)
VALUES (
    '5eed000000000000003001',
    '5eed000000000000002601',
    'Westphalian flax spinning',
    2000,
    48000,
    '1867-01-01 00:00:00.000000'
);

-- +goose Down
DELETE FROM domestic_industries WHERE id IN ('5eed000000000000003001');
