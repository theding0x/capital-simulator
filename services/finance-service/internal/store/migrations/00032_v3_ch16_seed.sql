-- +goose Up
-- Vol. III Ch. 16 seed — three commercial-capital exemplars.
-- Seed IDs: 5eed000000000000001601–5eed000000000000001603 (22-char, never collide with 24-char NewID()).
-- function: 0=realisation, 1=storage, 2=transport, 3=distribution

INSERT INTO commercial_capitals
    (id, money_advanced, commodity_description, `function`, surplus_value_produced, created_at)
VALUES
    ('5eed000000000000001601', 300000, 'Linen merchant: £3,000 advanced, 30,000 yards at 10p/yard', 0, 0, '2026-05-01 00:00:00.000000'),
    ('5eed000000000000001602', 600000, 'Industrial handoff: £600 commodity-capital (400c+100v+100s) realised by merchant', 0, 0, '2026-05-01 00:00:01.000000'),
    ('5eed000000000000001603', 150000, 'Storage variant: £1,500 stock held by merchant awaiting favourable price', 1, 0, '2026-05-01 00:00:02.000000');

-- +goose Down
DELETE FROM commercial_capitals WHERE id IN (
    '5eed000000000000001601',
    '5eed000000000000001602',
    '5eed000000000000001603'
);
