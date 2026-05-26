-- Vol. II Ch. 10 — Theories of Fixed and Circulating Capital:
-- The Physiocrats and Adam Smith.
--
-- Static reference tables for pre-Marxian economist terminology.
-- These tables are never written at runtime; all rows come from the seed migration.

-- +goose Up
CREATE TABLE IF NOT EXISTS economist_attributions (
    id           VARCHAR(24)  NOT NULL,
    concept      VARCHAR(255) NOT NULL,
    theorist     VARCHAR(100) NOT NULL,
    edition_year BIGINT       NOT NULL,
    anticipates  VARCHAR(100) NOT NULL,
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS economist_attribution_errors (
    attribution_id VARCHAR(24)  NOT NULL,
    error_code     VARCHAR(100) NOT NULL,
    PRIMARY KEY (attribution_id, error_code),
    CONSTRAINT fk_eae_attribution
        FOREIGN KEY (attribution_id)
        REFERENCES economist_attributions(id)
        ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS economist_attribution_errors;
DROP TABLE IF EXISTS economist_attributions;
