-- +goose Up
CREATE TABLE merchant_turnovers (
    id                    VARCHAR(24) NOT NULL,
    money_advanced        BIGINT NOT NULL,
    general_rate_bp       BIGINT NOT NULL,
    turnover_count        BIGINT NOT NULL,
    units_per_turnover    BIGINT NOT NULL,
    annual_merchant_profit BIGINT NOT NULL,
    markup_per_unit       BIGINT NOT NULL,
    created_at            DATETIME(6) NOT NULL,
    PRIMARY KEY (id),
    INDEX idx_merchant_turnovers_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS merchant_turnovers;
