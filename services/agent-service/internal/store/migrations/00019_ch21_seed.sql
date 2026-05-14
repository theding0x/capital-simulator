-- +goose Up
-- Seed data for Capital Vol. I, Ch. 21 — Piece-Wages.
--
-- Core Marx fixture (Ch. 21): "ordinary working-day 12 hours ... 24 pieces;
-- value of 24 pieces = 6s.; labourer receives 1½d. per piece; earns 3s. in 12 hours"
--
-- Using farthings (¼d.) as the unit:
--   daily wage        = 3s.  = 144 farthings
--   day value product = 6s.  = 288 farthings
--   normal output     = 24 pieces
--   piece price       = 144/24 = 6 farthings (1½d.)
--   piece value       = 288/24 = 12 farthings (3d.)
--
-- Sub-contract (sweating system):
--   Thomas Hobson (head) contracts at 6 farthings/piece,
--   pays assistants 4 farthings/piece — keeps spread of 2 farthings.

INSERT INTO piece_wages
    (id, agent_id, price_pence, normal_output, created_at)
VALUES
    ('5eed000000000000002100f1',
     '5eed000000000000001900ff',
     6, 24,
     '2026-05-14 09:00:00.000000');

INSERT INTO sub_contracts
    (id, head_labourer_id, assistant_ids, piece_rate_pence, assistant_rate_pence, created_at)
VALUES
    ('5eed000000000000002100f2',
     '5eed000000000000001900ff',
     '["5eed000000000000001900fe", "5eed000000000000001900fd"]',
     6, 4,
     '2026-05-14 09:01:00.000000');

-- +goose Down
DELETE FROM sub_contracts
WHERE id IN ('5eed000000000000002100f2');

DELETE FROM piece_wages
WHERE id IN ('5eed000000000000002100f1');
