-- Capital-simulator seed data.
-- Marx's canonical commodity examples from Capital, Vol. I, Ch. 1–3.
-- Mounted into /docker-entrypoint-initdb.d/ as seed.sql by docker-compose.
-- Runs once on first container boot, after init.sql creates the databases.
-- SNLT values are in minutes per unit.

-- ─── commodity database ──────────────────────────────────────────────────────

USE commodity;

CREATE TABLE IF NOT EXISTS commodities (
    id               VARCHAR(24)  NOT NULL PRIMARY KEY,
    name             VARCHAR(255) NOT NULL,
    use_value_desc   TEXT         NOT NULL,
    use_value_unit   VARCHAR(100) NOT NULL,
    cl_kind          VARCHAR(100) NOT NULL DEFAULT '',
    cl_description   TEXT         NOT NULL,
    snlt_per_unit    BIGINT       NOT NULL DEFAULT 0,
    created_at       DATETIME(6)  NOT NULL,
    updated_at       DATETIME(6)  NOT NULL,
    UNIQUE KEY uq_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Ch.1: 20 yards linen = 1 coat, both embodying 600 min SNLT (§1.1).
-- Ch.1 expanded form: iron (120 min/cwt), corn (240 min/qtr).
-- Ch.3: gold crystallises as the universal equivalent / money-commodity.
INSERT IGNORE INTO commodities
    (id, name, use_value_desc, use_value_unit, cl_kind, cl_description, snlt_per_unit, created_at, updated_at)
VALUES
    ('000000000000000000000001', 'linen', 'linen for clothing',     'yards', 'weaving',     '', 30,  '2024-01-01 00:00:00.000000', '2024-01-01 00:00:00.000000'),
    ('000000000000000000000002', 'coat',  'a coat to wear',          'coats', 'tailoring',   '', 600, '2024-01-01 00:00:00.000000', '2024-01-01 00:00:00.000000'),
    ('000000000000000000000003', 'iron',  'pig iron for tools',      'cwt',   'smelting',    '', 120, '2024-01-01 00:00:00.000000', '2024-01-01 00:00:00.000000'),
    ('000000000000000000000004', 'corn',  'corn for food',           'qtr',   'agriculture', '', 240, '2024-01-01 00:00:00.000000', '2024-01-01 00:00:00.000000'),
    ('000000000000000000000005', 'gold',  'gold as money-commodity', 'oz',    'mining',      '', 60,  '2024-01-01 00:00:00.000000', '2024-01-01 00:00:00.000000');

-- ─── market database ──────────────────────────────────────────────────────────

USE market;

CREATE TABLE IF NOT EXISTS owners (
    id         VARCHAR(24)  NOT NULL PRIMARY KEY,
    name       VARCHAR(255) NOT NULL,
    created_at DATETIME(6)  NOT NULL,
    updated_at DATETIME(6)  NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS offers (
    id                 VARCHAR(24)  NOT NULL PRIMARY KEY,
    owner_id           VARCHAR(24)  NOT NULL,
    commodity_id       VARCHAR(24)  NOT NULL,
    quantity           DOUBLE       NOT NULL,
    seeks_kind         VARCHAR(20)  NOT NULL DEFAULT '',
    seeks_commodity_id VARCHAR(24)  NOT NULL DEFAULT '',
    created_at         DATETIME(6)  NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS exchanges (
    id                    VARCHAR(24) NOT NULL PRIMARY KEY,
    giver_id              VARCHAR(24) NOT NULL,
    receiver_id           VARCHAR(24) NOT NULL,
    giver_commodity_id    VARCHAR(24) NOT NULL,
    giver_qty             DOUBLE      NOT NULL,
    receiver_commodity_id VARCHAR(24) NOT NULL,
    receiver_qty          DOUBLE      NOT NULL,
    realised_value        BIGINT      NOT NULL,
    created_at            DATETIME(6) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS market_config (
    key_name     VARCHAR(50) NOT NULL PRIMARY KEY,
    commodity_id VARCHAR(24) NOT NULL DEFAULT '',
    ts           DATETIME(6) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS prices (
    commodity_id       VARCHAR(24) NOT NULL PRIMARY KEY,
    money_commodity_id VARCHAR(24) NOT NULL,
    amount             BIGINT      NOT NULL,
    updated_at         DATETIME(6) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Commodity owners: weavers produce linen, tailors produce coats (Capital §1.4).
INSERT IGNORE INTO owners (id, name, created_at, updated_at)
VALUES
    ('000000000000000000000011', 'weaver', '2024-01-01 00:00:00.000000', '2024-01-01 00:00:00.000000'),
    ('000000000000000000000012', 'tailor', '2024-01-01 00:00:00.000000', '2024-01-01 00:00:00.000000');

-- Open barter offers: 20 yards linen for 1 coat (Capital §1.3).
INSERT IGNORE INTO offers
    (id, owner_id, commodity_id, quantity, seeks_kind, seeks_commodity_id, created_at)
VALUES
    ('000000000000000000000021', '000000000000000000000011', '000000000000000000000001', 20, 'coat',  '000000000000000000000002', '2024-01-01 00:00:00.000000'),
    ('000000000000000000000022', '000000000000000000000012', '000000000000000000000002',  1, 'linen', '000000000000000000000001', '2024-01-01 00:00:00.000000');

-- Realised exchange: 20 yards linen for 1 coat; 600 min realised value (Capital §1.1).
INSERT IGNORE INTO exchanges
    (id, giver_id, receiver_id, giver_commodity_id, giver_qty, receiver_commodity_id, receiver_qty, realised_value, created_at)
VALUES
    ('000000000000000000000031',
     '000000000000000000000011',
     '000000000000000000000012',
     '000000000000000000000001', 20,
     '000000000000000000000002',  1,
     600,
     '2024-01-01 00:00:00.000000');

-- Gold as universal equivalent and money-commodity (Capital §1.3, Ch.2–3).
INSERT IGNORE INTO market_config (key_name, commodity_id, ts)
VALUES
    ('universal_equivalent', '000000000000000000000005', '2024-01-01 00:00:00.000000'),
    ('money_commodity',      '000000000000000000000005', '2024-01-01 00:00:00.000000');

-- Prices in gold. amount = floor(commodity_snlt / gold_snlt) per unit.
-- coat: 600/60 = 10 oz gold. iron: 120/60 = 2 oz. corn: 240/60 = 4 oz.
-- linen: 30/60 = 0 (truncates at 1 unit; needs qty >= 2 for a non-zero price).
INSERT IGNORE INTO prices (commodity_id, money_commodity_id, amount, updated_at)
VALUES
    ('000000000000000000000002', '000000000000000000000005', 10, '2024-01-01 00:00:00.000000'),
    ('000000000000000000000003', '000000000000000000000005',  2, '2024-01-01 00:00:00.000000'),
    ('000000000000000000000004', '000000000000000000000005',  4, '2024-01-01 00:00:00.000000');
