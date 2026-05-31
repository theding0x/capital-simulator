-- +goose Up
CREATE TABLE gold_reserves (
    id          VARCHAR(24)  NOT NULL,
    name        VARCHAR(255) NOT NULL DEFAULT '',
    amount      BIGINT       NOT NULL DEFAULT 0,
    description TEXT         NOT NULL,
    created_at  DATETIME(6)  NOT NULL,
    PRIMARY KEY (id),
    INDEX idx_gold_reserves_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE rates_of_exchange (
    id           VARCHAR(24)  NOT NULL,
    name         VARCHAR(255) NOT NULL DEFAULT '',
    deviation_bp BIGINT       NOT NULL DEFAULT 0,
    description  TEXT         NOT NULL,
    created_at   DATETIME(6)  NOT NULL,
    PRIMARY KEY (id),
    INDEX idx_rates_of_exchange_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS rates_of_exchange;
DROP TABLE IF EXISTS gold_reserves;
