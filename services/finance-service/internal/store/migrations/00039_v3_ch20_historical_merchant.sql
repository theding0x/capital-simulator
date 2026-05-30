-- +goose Up
-- Vol. III Ch. 20 — Historical Facts about Merchant's Capital
-- stage: 1=pre_capitalist, 2=transition, 3=subordinated
-- subordination_index: basis points [0,10000] — merchant capital ÷ total social capital × 10000
-- wage_labour: TINYINT(1) — 0=false, 1=true; pre-capitalist rows must have wage_labour=0 (invariant 1)
CREATE TABLE historical_merchant_capitals (
    id                  VARCHAR(24)  NOT NULL,
    name                VARCHAR(255) NOT NULL,
    stage               TINYINT      NOT NULL,
    profit_source       VARCHAR(60)  NOT NULL,
    wage_labour         TINYINT(1)   NOT NULL DEFAULT 0,
    subordination_index BIGINT       NOT NULL,
    created_at          DATETIME(6)  NOT NULL,
    PRIMARY KEY (id),
    INDEX idx_hmc_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS historical_merchant_capitals;
