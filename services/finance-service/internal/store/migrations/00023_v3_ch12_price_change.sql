-- +goose Up
-- Vol. III Ch. 12 — Supplementary Remarks (on prices of production).
-- One table: price_of_production_changes. A commodity's price of production
-- changes for one of only two reasons — the general rate of profit changes, or
-- the commodity's own value changes — and the boolean flags record which.
-- Prices stored in whole units (optional context).

CREATE TABLE price_of_production_changes (
    id            VARCHAR(24)  NOT NULL,
    sphere_name   VARCHAR(120) NOT NULL,
    cause         VARCHAR(32)  NOT NULL, -- 'general_rate_change' | 'individual_value_change' | 'both' | ''
    rate_changed  TINYINT(1)   NOT NULL, -- general rate of profit changed
    value_changed TINYINT(1)   NOT NULL, -- commodity's own value changed
    price_changed TINYINT(1)   NOT NULL, -- price of production moved (cause != '')
    old_price     BIGINT       NOT NULL, -- whole units
    new_price     BIGINT       NOT NULL, -- whole units
    created_at    DATETIME(6)  NOT NULL,
    PRIMARY KEY (id),
    INDEX idx_price_of_production_changes_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS price_of_production_changes;
