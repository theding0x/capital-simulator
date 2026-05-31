-- +goose Up
-- Vol. III Ch. 27 seeds — Marx's credit-role exemplars.
-- Seed IDs: 5eed000000000000002700..2701 for stock_companies,
--           5eed000000000000002710..2711 for cooperative_factories.
--
-- Stock company 2700: United Alkali Trust (Marx's fixture, Ch. 27).
--   Fixed capital £5,000,000 = 500,000,000 pence.
--   Floating capital £1,000,000 = 100,000,000 pence.
--   Total capital £6,000,000 = 600,000,000 pence.
--   Manager: salaried director, not a profit-drawing owner.
--
-- Stock company 2701: Generic Railway Company — large fixed capital, small floating.
--   Fixed £9,000,000 (900,000,000 p); floating £1,000,000 (100,000,000 p); total £10,000,000.
--
-- Cooperative 2710: Rochdale Pioneers (1844) — 28 worker-owners, £28 capital (2800 p).
-- Cooperative 2711: Paisley Thread Works — larger worker co-op.

INSERT INTO `stock_companies`
    (`id`, `name`, `fixed_capital`, `floating_capital`, `total_capital`,
     `manager_name`, `manager_title`, `manager_wage_minutes`, `description`, `created_at`)
VALUES
    ('5eed000000000000002700', 'United Alkali Trust',    500000000, 100000000, 600000000,
     'Salaried Director', 'Managing Director', 480,
     'Joint-stock amalgamation of alkali firms; owner becomes pure money-capitalist; manager earns wages of superintendence',
     '2025-01-01 00:00:00.000000'),
    ('5eed000000000000002701', 'South Eastern Railway',  900000000, 100000000, 1000000000,
     'Railway Manager', 'General Manager', 600,
     'Large fixed capital in permanent way and rolling stock; functioning capitalist separated from ownership',
     '2025-01-01 00:01:00.000000');

INSERT INTO `cooperative_factories`
    (`id`, `name`, `worker_owned`, `worker_count`, `capital_advanced`, `description`, `created_at`)
VALUES
    ('5eed000000000000002710', 'Rochdale Pioneers Cooperative', 1, 28, 2800,
     'First modern workers cooperative (1844); positive sprout of the new mode of production within the old',
     '2025-01-01 00:00:00.000000'),
    ('5eed000000000000002711', 'Paisley Thread Works Cooperative', 1, 120, 500000,
     'Worker-owned thread mill; associated producers controlling their own conditions of production',
     '2025-01-01 00:01:00.000000');

-- +goose Down
DELETE FROM `cooperative_factories`
WHERE `id` IN ('5eed000000000000002710', '5eed000000000000002711');

DELETE FROM `stock_companies`
WHERE `id` IN ('5eed000000000000002700', '5eed000000000000002701');
