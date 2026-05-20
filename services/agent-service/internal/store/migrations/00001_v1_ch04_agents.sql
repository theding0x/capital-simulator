-- +goose Up
CREATE TABLE IF NOT EXISTS agents (
    id            VARCHAR(24)  NOT NULL PRIMARY KEY,
    name          VARCHAR(255) NOT NULL,
    class         VARCHAR(50)  NOT NULL,
    money_balance BIGINT       NOT NULL DEFAULT 0,
    hoarding      TINYINT(1)   NOT NULL DEFAULT 0,
    created_at    DATETIME(6)  NOT NULL,
    updated_at    DATETIME(6)  NOT NULL,
    INDEX idx_class (class)
);

-- +goose Down
DROP TABLE IF EXISTS agents;
