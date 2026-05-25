-- +goose Up
-- Vol. II Ch. 6 seed — Marx's fixtures for dashboard population.
-- Seed IDs: 5eed000000000000 + 06 (chapter) + sequence.

-- §I(a): Spinning Mill 1871 seller — 480 necessary + 120 surplus minutes.
INSERT INTO circulation_agents (id, industrial_capital_id, agent_id, role, wage_rate_pence, labour_minutes_necessary, labour_minutes_surplus)
VALUES ('5eed000000000000060101', '5eed000000000000000401', '', 'seller', 480, 480, 120);

-- §I(b): Spinning Mill 1871 book-keeper — 1,000 labour + 200 stationery = 1,200 pence/week.
INSERT INTO circulation_costs (id, industrial_capital_id, kind, nature, pence, incurred_at)
VALUES ('5eed000000000000060201', '5eed000000000000000401', 'book_keeping', 'value_preserving', 1200, '1871-01-01 00:00:00.000');

-- §I(c): National gold reserve — 1,000,000 pence reserve, 10,000 pence annual replacement (1% wear).
INSERT INTO money_as_faux_frais (id, industrial_capital_id, pence, annual_replacement_pence)
VALUES ('5eed000000000000060301', '', 1000000, 10000);

-- §II(b): Voluntary commodity supply — 100,000 pence of yarn held for steady sale.
INSERT INTO commodity_supplies (id, industrial_capital_id, commodity_id, pence, held_since, is_voluntary)
VALUES ('5eed000000000000060401', '5eed000000000000000401', '', 100000, '1871-01-01 00:00:00.000', 1);

-- Storage cost for voluntary supply: building 500 + labour 300 + preservation 200 = 1,000.
INSERT INTO storage_costs (id, commodity_supply_id, building_pence, labour_pence, preservation_labour_pence)
VALUES ('5eed000000000000060501', '5eed000000000000060401', 500, 300, 200);

-- §II(b): Involuntary supply (stagnation) — identical pence, is_voluntary = 0.
INSERT INTO commodity_supplies (id, industrial_capital_id, commodity_id, pence, held_since, is_voluntary)
VALUES ('5eed000000000000060601', '5eed000000000000000401', '', 100000, '1871-01-01 00:00:00.000', 0);

-- Storage cost for involuntary supply — same amounts, but a pure loss.
INSERT INTO storage_costs (id, commodity_supply_id, building_pence, labour_pence, preservation_labour_pence)
VALUES ('5eed000000000000060701', '5eed000000000000060601', 500, 300, 200);

-- §III Birmingham glass (historic): base 24 pence/ton-mile, no breakage multiplier.
INSERT INTO transport_tariffs (id, commodity_id, base_pence_per_ton_mile, fragility_multiplier_basis_points, breakage_risk_multiplier_basis_points)
VALUES ('5eed000000000000060801', 'glass-birmingham-historic', 24, 0, 10000);

-- §III Birmingham glass (modern): same base rate, 3× breakage risk (30,000 bp).
INSERT INTO transport_tariffs (id, commodity_id, base_pence_per_ton_mile, fragility_multiplier_basis_points, breakage_risk_multiplier_basis_points)
VALUES ('5eed000000000000060901', 'glass-birmingham-modern', 24, 0, 30000);

-- Historic transport leg: £11/crate glass, 50 miles Birmingham→London.
INSERT INTO transport_legs (id, commodity_id, quantity, origin_id, destination_id, distance_meters, weight_grams, volume_cubic_millimeters, labour_cost_pence, means_of_transport_pence, surplus_pence)
VALUES ('5eed000000000000060a01', 'glass-birmingham-historic', 1, 'birmingham', 'london', 80450, 1000000, 0, 100, 20, 0);

-- Modern transport leg: £2/crate glass, same route — 3× tariff (breakage risk).
INSERT INTO transport_legs (id, commodity_id, quantity, origin_id, destination_id, distance_meters, weight_grams, volume_cubic_millimeters, labour_cost_pence, means_of_transport_pence, surplus_pence)
VALUES ('5eed000000000000060b01', 'glass-birmingham-modern', 1, 'birmingham', 'london', 80450, 1000000, 0, 300, 60, 0);

-- +goose Down
DELETE FROM transport_legs WHERE id IN ('5eed000000000000060a01','5eed000000000000060b01');
DELETE FROM transport_tariffs WHERE id IN ('5eed000000000000060801','5eed000000000000060901');
DELETE FROM storage_costs WHERE id IN ('5eed000000000000060501','5eed000000000000060701');
DELETE FROM commodity_supplies WHERE id IN ('5eed000000000000060401','5eed000000000000060601');
DELETE FROM money_as_faux_frais WHERE id = '5eed000000000000060301';
DELETE FROM circulation_costs WHERE id = '5eed000000000000060201';
DELETE FROM circulation_agents WHERE id = '5eed000000000000060101';
