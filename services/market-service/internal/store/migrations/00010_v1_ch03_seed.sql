-- +goose Up
-- Vol. I Ch. 3 seed — Marx-faithful exemplars so the Ch. 3 panel comes up
-- populated on a fresh MySQL volume.
--
-- Seed ID convention: 5eed00000000000000<CC><XX> where CC=03 (chapter) and
-- XX is the row sequence. These IDs never collide with market.NewHoardID()
-- output (96-bit hex from crypto/rand).
--
-- These seeds reference owner IDs that the agent / market seed migrations
-- create. If the referenced owners do not yet exist on a fresh install,
-- the INSERT will still succeed — the columns are not declared as foreign
-- keys, mirroring the existing market-service tables.

INSERT INTO hoards (id, owner_id, amount, created_at) VALUES
    ('5eed000000000000000301', '5eed000000000000000202', 12000,
     '2026-05-01 00:00:00.000000'),
    ('5eed000000000000000302', '5eed000000000000000202', 48000,
     '2026-05-15 00:00:00.000000');

INSERT INTO payment_obligations (id, creditor_id, debtor_id, amount, created_at, due_at, paid_at) VALUES
    ('5eed000000000000000311', '5eed000000000000000201', '5eed000000000000000203', 24000,
     '2026-05-01 00:00:00.000000', '2026-05-30 00:00:00.000000', NULL),
    ('5eed000000000000000312', '5eed000000000000000201', '5eed000000000000000204', 8000,
     '2026-04-15 00:00:00.000000', '2026-05-15 00:00:00.000000',
     '2026-05-14 00:00:00.000000');

INSERT INTO world_money_transfers (id, sender_id, receiver_id, gold_mg, created_at) VALUES
    ('5eed000000000000000321', '5eed000000000000000205', '5eed000000000000000201', 15500,
     '2026-04-01 00:00:00.000000'),
    ('5eed000000000000000322', '5eed000000000000000201', '5eed000000000000000206', 9800,
     '2026-04-20 00:00:00.000000');

-- +goose Down
DELETE FROM world_money_transfers WHERE id IN (
    '5eed000000000000000321',
    '5eed000000000000000322'
);
DELETE FROM payment_obligations WHERE id IN (
    '5eed000000000000000311',
    '5eed000000000000000312'
);
DELETE FROM hoards WHERE id IN (
    '5eed000000000000000301',
    '5eed000000000000000302'
);
