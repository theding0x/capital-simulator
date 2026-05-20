-- +goose Up
-- Seed data for Capital Vol. I, Ch. 4 — The General Formula for Capital.
--
-- Three class positions are represented:
--   Capitalists (James Watt, Robert Peel) — advance money into commodities and
--     extract surplus-value, so their balances already include ∆M from the
--     circuits recorded below.
--   Workers (John, Ann) — sell a commodity to obtain money and immediately
--     spend it; no surplus emerges (C—M—C).
--   Miser (Thomas) — hoards gold, withdrawing money from circulation.
--
-- All money values are in pence (100 pence = £1.00).
-- Commodity IDs reference the commodity service (seeded separately):
--   000000000000000000000001 = linen
--   000000000000000000000002 = coat
--   000000000000000000000004 = corn

INSERT INTO agents (id, name, class, money_balance, hoarding, created_at, updated_at) VALUES
  -- Capitalists: balances reflect surplus already realised in circuits below.
  --   James Watt: started £500, two linen circuits (+£10, +£11) → £521.00
  ('5eed00000000000000000001', 'James Watt',  'capitalist', 52100,  0, '2026-05-01 09:00:00.000000', '2026-05-01 11:00:00.000000'),
  --   Robert Peel: started £250, one corn circuit (+£20) → £270.00
  ('5eed00000000000000000002', 'Robert Peel', 'capitalist', 27000,  0, '2026-05-01 09:00:00.000000', '2026-05-01 10:30:00.000000'),
  -- Workers: modest balances, no capital accumulation.
  ('5eed00000000000000000003', 'John',        'worker',      1200,  0, '2026-05-01 09:00:00.000000', '2026-05-01 09:00:00.000000'),
  ('5eed00000000000000000004', 'Ann',         'worker',       800,  0, '2026-05-01 09:00:00.000000', '2026-05-01 09:00:00.000000'),
  -- Miser: large gold hoard, withdrawn from circulation.
  ('5eed00000000000000000005', 'Thomas',      'miser',     100000,  1, '2026-05-01 09:00:00.000000', '2026-05-01 09:00:00.000000');

INSERT INTO capital_circuits (id, agent_id, m_advanced, commodity_id, m_returned, surplus_value, circuit_type, created_at) VALUES
  -- James Watt: M—C—M′, first circuit. £100 → linen → £110 (∆M = £10).
  ('5eed0000000000000000a101', '5eed00000000000000000001', 10000, '000000000000000000000001', 11000, 1000, 'M-C-M-prime', '2026-05-01 10:00:00.000000'),
  -- James Watt: M—C—M′, reinvested. £110 → linen → £121 (∆M = £11).
  ('5eed0000000000000000a102', '5eed00000000000000000001', 11000, '000000000000000000000001', 12100, 1100, 'M-C-M-prime', '2026-05-01 11:00:00.000000'),
  -- Robert Peel: M—C—M′. £200 → corn → £220 (∆M = £20).
  ('5eed0000000000000000b101', '5eed00000000000000000002', 20000, '000000000000000000000004', 22000, 2000, 'M-C-M-prime', '2026-05-01 10:30:00.000000'),
  -- John (worker): C—M—C. Sells coat (£5) → buys corn (£5). No surplus.
  ('5eed0000000000000000c101', '5eed00000000000000000003',   500, '000000000000000000000002',   500,    0, 'C-M-C',       '2026-05-01 10:15:00.000000');

-- +goose Down
DELETE FROM capital_circuits WHERE id IN (
  '5eed0000000000000000a101',
  '5eed0000000000000000a102',
  '5eed0000000000000000b101',
  '5eed0000000000000000c101'
);
DELETE FROM agents WHERE id IN (
  '5eed00000000000000000001',
  '5eed00000000000000000002',
  '5eed00000000000000000003',
  '5eed00000000000000000004',
  '5eed00000000000000000005'
);
