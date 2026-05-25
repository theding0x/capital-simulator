-- +goose Up
-- Vol. II Ch. 8 seed — Marx's §I fixtures for fixed/circulating capital.
--
-- Fixture 1: "If a machine worth £10,000 lasts for ten years..."
--   CapitalComponent (instrument_of_labour / fixed), pence = 1,000,000
--   FixedCapitalItem: service_life_minutes = 10 × 525,600 = 5,256,000
--
-- Fixture 2: British railway subcomponents of unequal durability —
--   rails (2yr), sleepers (8yr), earthworks (80yr).
--
-- Fixture 3: Auxiliary material (lubricating oil) and raw material (iron ore) —
--   both circulating.

-- Seed IDs: 5eed000000000000 08 XX

-- CapitalComponents
INSERT INTO capital_components (id, industrial_capital_id, kind, role, pence, entered_at) VALUES
('5eed00000000000008c1', '', 'fixed',       'instrument_of_labour', 1000000, '1867-01-01 00:00:00.000'),
('5eed00000000000008c2', '', 'circulating', 'raw_material',          50000,  '1867-01-01 00:00:00.000'),
('5eed00000000000008c3', '', 'circulating', 'auxiliary_material',    5000,   '1867-01-01 00:00:00.000'),
('5eed00000000000008c4', '', 'fixed',       'instrument_of_labour',  400000, '1867-01-01 00:00:00.000');

-- FixedCapitalItems
-- Fixture 1: £10,000 machine (10 years = 5,256,000 minutes)
INSERT INTO fixed_capital_items (id, capital_component_id, description, pence_purchased, service_life_minutes, wear_model, purchased_at) VALUES
('5eed00000000000008f1', '5eed00000000000008c1', 'Ten-year spinning machine (Marx §I fixture)', 1000000, 5256000, 'linear', '1867-01-01 00:00:00.000');

-- Fixture 2: Railway rolling stock (£4,000; 8-year life = 4,204,800 minutes)
INSERT INTO fixed_capital_items (id, capital_component_id, description, pence_purchased, service_life_minutes, wear_model, purchased_at) VALUES
('5eed00000000000008f2', '5eed00000000000008c4', 'Railway rolling stock — British railway fixture', 400000, 4204800, 'linear', '1867-01-01 00:00:00.000');

-- SinkingFundsForItems (target = pence_purchased; accumulated starts at 0)
INSERT INTO sinking_funds_for_items (id, fixed_capital_item_id, accumulated_pence, target_pence) VALUES
('5eed00000000000008s1', '5eed00000000000008f1', 0, 1000000),
('5eed00000000000008s2', '5eed00000000000008f2', 0, 400000);

-- FixedCapitalSubcomponents (railway — unequal durability)
-- Rails: ~2-year life (1,051,200 min); sleepers: 8yr (4,204,800 min); earthworks: 80yr (42,048,000 min)
INSERT INTO fixed_capital_subcomponents (id, fixed_capital_item_id, name, pence, service_life_minutes, wear_model, last_replaced_at) VALUES
('5eed00000000000008u1', '5eed00000000000008f2', 'Rails',      100000, 1051200,  'linear', '1867-01-01 00:00:00.000'),
('5eed00000000000008u2', '5eed00000000000008f2', 'Sleepers',   150000, 4204800,  'linear', '1867-01-01 00:00:00.000'),
('5eed00000000000008u3', '5eed00000000000008f2', 'Earthworks', 150000, 42048000, 'linear', '1867-01-01 00:00:00.000');

-- +goose Down
DELETE FROM fixed_capital_subcomponents WHERE id IN (
    '5eed00000000000008u1', '5eed00000000000008u2', '5eed00000000000008u3'
);
DELETE FROM sinking_funds_for_items WHERE id IN (
    '5eed00000000000008s1', '5eed00000000000008s2'
);
DELETE FROM fixed_capital_items WHERE id IN (
    '5eed00000000000008f1', '5eed00000000000008f2'
);
DELETE FROM capital_components WHERE id IN (
    '5eed00000000000008c1', '5eed00000000000008c2',
    '5eed00000000000008c3', '5eed00000000000008c4'
);
