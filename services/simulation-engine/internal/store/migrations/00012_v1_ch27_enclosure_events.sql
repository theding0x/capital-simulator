-- +goose Up
CREATE TABLE enclosure_events (
    id                   CHAR(24)     NOT NULL PRIMARY KEY,
    period               VARCHAR(255) NOT NULL,
    acres_enclosed       BIGINT       NOT NULL,
    population_displaced BIGINT       NOT NULL,
    beneficiary          VARCHAR(255) NOT NULL,
    created_at           DATETIME(6)  NOT NULL,
    KEY idx_enclosure_events_created (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- +goose Down
DROP TABLE IF EXISTS enclosure_events;