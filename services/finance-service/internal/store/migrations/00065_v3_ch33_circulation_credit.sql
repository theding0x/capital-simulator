-- +goose Up
CREATE TABLE clearing_house_settlements (
    id          VARCHAR(24)  NOT NULL,
    name        VARCHAR(255) NOT NULL DEFAULT '',
    description TEXT         NOT NULL,
    total_claims BIGINT      NOT NULL DEFAULT 0,
    money_used  BIGINT       NOT NULL DEFAULT 0,
    net_balance BIGINT       NOT NULL DEFAULT 0,
    created_at  DATETIME(6)  NOT NULL,
    PRIMARY KEY (id),
    INDEX idx_clearing_house_settlements_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS clearing_house_settlements;
