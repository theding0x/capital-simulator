-- +goose Up
-- Vol. III Ch. 16 — Commercial Capital. Opens Part IV.

CREATE TABLE commercial_capitals (
    id                     VARCHAR(24)   NOT NULL,
    money_advanced         BIGINT        NOT NULL,
    commodity_description  VARCHAR(255)  NOT NULL DEFAULT '',
    `function`             TINYINT       NOT NULL DEFAULT 0,
    surplus_value_produced BIGINT        NOT NULL DEFAULT 0,
    created_at             DATETIME(6)   NOT NULL,
    PRIMARY KEY (id),
    INDEX idx_commercial_capitals_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS commercial_capitals;
