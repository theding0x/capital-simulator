-- +goose Up
-- Vol. III Ch. 17 — Commercial Profit. The merchant advances productive capital
-- (borrowed from the industrial sphere) and commercial capital; the general rate
-- of profit is computed over both. NewValueCreated is always 0.

CREATE TABLE commercial_profits (
    id                  VARCHAR(24)  NOT NULL,
    productive_capital  BIGINT       NOT NULL,
    commercial_capital  BIGINT       NOT NULL,
    surplus_value       BIGINT       NOT NULL,
    unadjusted_rate_bp  BIGINT       NOT NULL,
    adjusted_rate_bp    BIGINT       NOT NULL,
    new_value_created   BIGINT       NOT NULL DEFAULT 0,
    created_at          DATETIME(6)  NOT NULL,
    PRIMARY KEY (id),
    INDEX idx_commercial_profits_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS commercial_profits;
