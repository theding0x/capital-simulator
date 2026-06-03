-- +goose Up
-- Vol. I Ch. 3 seed — Marx-faithful C—M—C circuits so the circuit recorder on
-- the Ch. 3 panel comes up populated on a fresh MySQL volume.
--
-- Seed ID convention: 5eed00000000000000<CC><XX> where CC=03 (chapter) and XX
-- is the row sequence. Circuit rows use 33XX; the commodity / money IDs the
-- legs reference use 34XX. These never collide with market.NewCircuitID()
-- output (96-bit hex from crypto/rand). Owner IDs reuse the seed owners the
-- existing Ch. 3 seed (00010) already references.
--
-- Both circuits honour the load-bearing invariant sale_value == purchase_value:
-- the same value-magnitude (£2 = 200 pence) leaves in commodity-form and returns
-- in commodity-form. The chain is Marx's own §2 example:
--   wheat — £2 — linen — £2 — Bible — £2 — brandy.

INSERT INTO circuits
    (id, sale_kind, sale_commodity_id, sale_money_id, sale_owner_id, sale_value,
     purchase_kind, purchase_commodity_id, purchase_money_id, purchase_owner_id, purchase_value, created_at)
VALUES
    -- The weaver sells 20 yards of linen for £2 and lays the £2 out on a family Bible.
    ('5eed000000000000000331',
     'sale',     '5eed000000000000000341', '5eed000000000000000343', '5eed000000000000000201', 200,
     'purchase', '5eed000000000000000342', '5eed000000000000000343', '5eed000000000000000201', 200,
     '2026-05-01 00:00:00.000000'),
    -- The farmer sells a quarter of wheat for £2 and buys 20 yards of linen.
    ('5eed000000000000000332',
     'sale',     '5eed000000000000000344', '5eed000000000000000343', '5eed000000000000000202', 200,
     'purchase', '5eed000000000000000341', '5eed000000000000000343', '5eed000000000000000202', 200,
     '2026-05-02 00:00:00.000000');

-- +goose Down
DELETE FROM circuits WHERE id IN (
    '5eed000000000000000331',
    '5eed000000000000000332'
);
