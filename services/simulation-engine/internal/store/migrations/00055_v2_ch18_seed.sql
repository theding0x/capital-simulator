-- +goose Up
-- Vol. II Ch. 18 — seed fixtures for the money-capital reproduction model.
--
-- Marx's two-department scheme (simplified, 1865 illustrative period):
--   Total social capital = £1,000 (10,000,000 pence, 1 GBP = 240 pence)
--   Dept I  (means of production):  40% of total reserve
--   Dept II (articles of consumption): 35% of total reserve
--   Wage-rotation fund: 20% — advanced before output is sold
--   Idle hoard: 5% — latent money capital not yet in circuit
--
-- Money velocity: 2× per year (20000 basis points), so £1 000 in stock
-- circulates as £2 000 of effective demand.
--
-- Wage-rotation funds: both departments pay weekly (52 cycles p.a.).
-- Inter-department flows: I(v+s) → II(c) and II(c) → I at £208.33 each.
--
-- Seed IDs follow the pattern 5eed000000000000018XX so they are
-- recognisable and never collide with NewXxxID() outputs.

INSERT INTO money_supply_apportionments
    (id, total_social_money_pence, dept_i_reserve_pence,
     dept_ii_reserve_pence, wage_rotation_pence, idle_hoard_pence, period)
VALUES
    ('5eed000000000000001801', 10000000, 4000000, 3500000, 2000000, 500000, '1865');

INSERT INTO department_money_reserves
    (id, department, reserve_pence, reserve_purpose, period)
VALUES
    ('5eed000000000000001802', 'department_i',  2000000, 'wage_payment',        '1865'),
    ('5eed000000000000001803', 'department_ii', 1500000, 'means_of_production', '1865');

INSERT INTO circulating_money_masses
    (id, money_stock_pence, velocity_per_year_basis_points,
     effective_circulating_value_pence, period)
VALUES
    ('5eed000000000000001804', 10000000, 20000, 20000000, '1865');

INSERT INTO wage_rotation_funds
    (id, fund_pence, wage_cycle_frequency, department, period)
VALUES
    ('5eed000000000000001805', 1000000, 52, 'department_i',  '1865'),
    ('5eed000000000000001806', 1000000, 52, 'department_ii', '1865');

INSERT INTO inter_department_settlements
    (id, from_department, to_department, settled_pence, settlement_purpose, period)
VALUES
    ('5eed000000000000001807', 'department_i',  'department_ii', 2000000, 'means_of_production', '1865'),
    ('5eed000000000000001808', 'department_ii', 'department_i',  2000000, 'means_of_production', '1865');

-- +goose Down
DELETE FROM inter_department_settlements
    WHERE id IN ('5eed000000000000001807', '5eed000000000000001808');
DELETE FROM wage_rotation_funds
    WHERE id IN ('5eed000000000000001805', '5eed000000000000001806');
DELETE FROM circulating_money_masses
    WHERE id = '5eed000000000000001804';
DELETE FROM department_money_reserves
    WHERE id IN ('5eed000000000000001802', '5eed000000000000001803');
DELETE FROM money_supply_apportionments
    WHERE id = '5eed000000000000001801';
