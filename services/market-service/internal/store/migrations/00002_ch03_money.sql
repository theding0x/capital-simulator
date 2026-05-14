-- +goose Up
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

-- +goose Down
DROP TABLE IF EXISTS prices;
DROP TABLE IF EXISTS market_config;
