-- +goose Up
CREATE TABLE dr1_tables (
    id               VARCHAR(24)  NOT NULL,
    description      VARCHAR(500) NOT NULL DEFAULT '',
    total_rental_bp  BIGINT       NOT NULL DEFAULT 0,
    regulating_grade TINYINT      NOT NULL DEFAULT 1,
    entries_json     MEDIUMTEXT   NOT NULL,
    created_at       DATETIME(6)  NOT NULL,
    PRIMARY KEY (id),
    INDEX idx_dr1_tables_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE location_rent_factors (
    id                 VARCHAR(24) NOT NULL,
    parcel_id          VARCHAR(24) NOT NULL DEFAULT '',
    transport_cost_bp  BIGINT      NOT NULL DEFAULT 0,
    rent_equivalent_bp BIGINT      NOT NULL DEFAULT 0,
    created_at         DATETIME(6) NOT NULL,
    PRIMARY KEY (id),
    INDEX idx_location_rent_factors_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS location_rent_factors;
DROP TABLE IF EXISTS dr1_tables;
