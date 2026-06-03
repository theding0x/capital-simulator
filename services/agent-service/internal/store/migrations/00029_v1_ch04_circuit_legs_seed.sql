-- +goose Up
-- Seed data for Capital Vol. I, Ch. 4 — coordinated multi-agent exchange (#216).
--
-- A single market exchange (20 yards of linen = 1 coat, realised at 1200
-- labour-minutes) lands two complementary legs: John (worker, seed id …0003)
-- sells the linen to Ann (worker, seed id …0004). John closes a sale leg
-- (C → M); Ann closes a purchase leg (M → C). Both reference the linen
-- commodity (000…001, seeded in commodity-service) and the same exchange id.
INSERT INTO circuit_legs (id, agent_id, kind, counterparty_id, commodity_id, value_minutes, exchange_id, created_at) VALUES
  ('5eed0000000000000000d101', '5eed00000000000000000003', 'sale',     '5eed00000000000000000004', '000000000000000000000001', 1200, '5eed0000000000000000e001', '2026-05-01 10:20:00.000000'),
  ('5eed0000000000000000d102', '5eed00000000000000000004', 'purchase', '5eed00000000000000000003', '000000000000000000000001', 1200, '5eed0000000000000000e001', '2026-05-01 10:20:00.000000');

-- +goose Down
DELETE FROM circuit_legs WHERE id IN (
  '5eed0000000000000000d101',
  '5eed0000000000000000d102'
);
