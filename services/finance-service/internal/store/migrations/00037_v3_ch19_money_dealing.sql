-- +goose Up
-- Vol. III Ch. 19 — Money-Dealing Capital
-- money_dealing_capitals stores MoneyDealingCapital persisted entities.
-- operation_kinds is a JSON-encoded array of MoneyOperationKind integers (1–5).
CREATE TABLE money_dealing_capitals (
    id                   VARCHAR(24)   NOT NULL,
    money_advanced       BIGINT        NOT NULL,
    general_rate_bp      BIGINT        NOT NULL,
    operation_kinds      MEDIUMTEXT    NOT NULL,
    money_dealing_profit BIGINT        NOT NULL,
    new_value_created    BIGINT        NOT NULL DEFAULT 0,
    created_at           DATETIME(6)   NOT NULL,
    PRIMARY KEY (id),
    INDEX idx_mdc_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS money_dealing_capitals;
