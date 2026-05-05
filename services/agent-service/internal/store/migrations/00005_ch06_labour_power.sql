-- +goose Up
CREATE TABLE IF NOT EXISTS labour_workers (
    id                       VARCHAR(24)  NOT NULL PRIMARY KEY,
    kind                     VARCHAR(50)  NOT NULL DEFAULT 'worker',
    owns_labour_power        TINYINT(1)   NOT NULL DEFAULT 1,
    owns_commodities_to_sell TINYINT(1)   NOT NULL DEFAULT 0,
    capacity_minutes_per_day BIGINT       NOT NULL DEFAULT 0,
    created_at               DATETIME(6)  NOT NULL,
    updated_at               DATETIME(6)  NOT NULL
);

CREATE TABLE IF NOT EXISTS labour_capitalists (
    id            VARCHAR(24)  NOT NULL PRIMARY KEY,
    kind          VARCHAR(50)  NOT NULL DEFAULT 'capitalist',
    money_capital BIGINT       NOT NULL DEFAULT 0,
    created_at    DATETIME(6)  NOT NULL,
    updated_at    DATETIME(6)  NOT NULL
);

CREATE TABLE IF NOT EXISTS labour_power_offerings (
    id                       VARCHAR(24)  NOT NULL PRIMARY KEY,
    owner_id                 VARCHAR(24)  NOT NULL,
    capacity_minutes_per_day BIGINT       NOT NULL,
    contract_days            BIGINT       NOT NULL,
    asking_wage              BIGINT       NOT NULL,
    created_at               DATETIME(6)  NOT NULL,
    INDEX idx_owner_id (owner_id)
);

CREATE TABLE IF NOT EXISTS labour_power_purchases (
    id            VARCHAR(24)  NOT NULL PRIMARY KEY,
    seller_id     VARCHAR(24)  NOT NULL,
    buyer_id      VARCHAR(24)  NOT NULL,
    wage_minutes  BIGINT       NOT NULL,
    contract_days BIGINT       NOT NULL,
    created_at    DATETIME(6)  NOT NULL,
    INDEX idx_seller_id (seller_id),
    INDEX idx_buyer_id  (buyer_id)
);

-- +goose Down
DROP TABLE IF EXISTS labour_power_purchases;
DROP TABLE IF EXISTS labour_power_offerings;
DROP TABLE IF EXISTS labour_capitalists;
DROP TABLE IF EXISTS labour_workers;
