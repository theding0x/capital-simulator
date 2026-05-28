-- +goose Up
-- Vol. II Ch. 20 — seed fixtures for the simple-reproduction scheme.
--
-- Marx's canonical fixture (Capital Vol. II Ch. 20, pence × 1 000):
--   Dept I:  4 000c + 1 000v + 1 000s = 6 000 total
--   Dept II: 2 000c +   500v +   500s = 3 000 total
--   Keystone: I(v+s) = 2 000 == II(c) = 2 000  ✓
--
-- Three exchange flows:
--   1. Intra-Dept-I replacement: I buys MP from itself (I.c = 4 000 000)
--   2. I(v+s) → II: workers/capitalists of Dept I buy consumption goods (2 000 000)
--   3. II(c) → I: Dept II buys MP from Dept I (2 000 000)
--
-- Seed IDs follow the pattern 5eed000000000000002001, 2002, …

INSERT INTO simple_reproduction_schemes (id, period, is_balanced, tick_count)
VALUES ('5eed000000000000002001', '1865', 1, 0);

INSERT INTO departmental_capitals
    (id, scheme_id, department, constant_pence, variable_pence, surplus_pence, total_pence)
VALUES
    ('5eed000000000000002002', '5eed000000000000002001', 'department_i',
     4000000, 1000000, 1000000, 6000000),
    ('5eed000000000000002003', '5eed000000000000002001', 'department_ii',
     2000000, 500000, 500000, 3000000);

INSERT INTO interdepartment_exchanges
    (id, scheme_id, from_department, to_department, pence, kind, description)
VALUES
    ('5eed000000000000002004', '5eed000000000000002001',
     'department_i', 'department_i', 4000000, 'goods_for_money',
     'Intra-Dept-I: Dept I replaces its own constant capital internally (I.c = 4 000)'),
    ('5eed000000000000002005', '5eed000000000000002001',
     'department_i', 'department_ii', 2000000, 'goods_for_money',
     'I(v+s) → II: Dept I ships consumption goods; money flows II→I'),
    ('5eed000000000000002006', '5eed000000000000002001',
     'department_ii', 'department_i', 2000000, 'money_for_goods',
     'II(c) → I: Dept II advances money; receives means of production from Dept I');

INSERT INTO money_closed_loops (id, scheme_id, period, is_closed, net_flow_pence)
VALUES ('5eed000000000000002007', '5eed000000000000002001', '1865', 1, 0);

-- +goose Down
DELETE FROM money_closed_loops
    WHERE id = '5eed000000000000002007';
DELETE FROM interdepartment_exchanges
    WHERE id IN (
        '5eed000000000000002004',
        '5eed000000000000002005',
        '5eed000000000000002006'
    );
DELETE FROM departmental_capitals
    WHERE id IN ('5eed000000000000002002', '5eed000000000000002003');
DELETE FROM simple_reproduction_schemes
    WHERE id = '5eed000000000000002001';
