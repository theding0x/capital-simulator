-- +goose Up
-- Vol. I Ch. 3 — the two-leg C—M—C circuit of simple circulation. Each row is
-- one circuit: a sale (C → M) followed by a purchase (M → C), with the two legs
-- flattened into sale_* / purchase_* columns. The defining invariant
-- (sale_value == purchase_value: the same value-magnitude leaves in
-- commodity-form and returns in commodity-form) is enforced in the domain
-- layer (market.Circuit.Validate).

CREATE TABLE IF NOT EXISTS circuits (
    id                    VARCHAR(24) NOT NULL PRIMARY KEY,
    sale_kind             VARCHAR(16) NOT NULL,
    sale_commodity_id     VARCHAR(24) NOT NULL,
    sale_money_id         VARCHAR(24) NOT NULL,
    sale_owner_id         VARCHAR(24) NOT NULL,
    sale_value            BIGINT      NOT NULL,
    purchase_kind         VARCHAR(16) NOT NULL,
    purchase_commodity_id VARCHAR(24) NOT NULL,
    purchase_money_id     VARCHAR(24) NOT NULL,
    purchase_owner_id     VARCHAR(24) NOT NULL,
    purchase_value        BIGINT      NOT NULL,
    created_at            DATETIME(6) NOT NULL,
    INDEX idx_circuits_created_at (created_at),
    INDEX idx_circuits_sale_owner (sale_owner_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS circuits;
