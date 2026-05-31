-- +goose Up
CREATE TABLE usurers_capitals (
    id               VARCHAR(24)  NOT NULL,
    name             VARCHAR(255) NOT NULL DEFAULT '',
    stage            BIGINT       NOT NULL DEFAULT 1,
    wage_labour      TINYINT(1)   NOT NULL DEFAULT 0,
    borrowers_json   TEXT         NOT NULL,
    interest_rate_bp BIGINT       NOT NULL DEFAULT 0,
    subordinated_to  VARCHAR(255) NOT NULL DEFAULT '',
    description      TEXT         NOT NULL,
    created_at       DATETIME(6)  NOT NULL,
    PRIMARY KEY (id),
    INDEX idx_usurers_capitals_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS usurers_capitals;
