-- +goose Up
CREATE TABLE land_parcels (
    id               VARCHAR(24)  NOT NULL,
    owner_id         VARCHAR(24)  NOT NULL DEFAULT '',
    fertility_grade  BIGINT       NOT NULL DEFAULT 1,
    area_hectares    BIGINT       NOT NULL DEFAULT 0,
    location         VARCHAR(255) NOT NULL DEFAULT '',
    created_at       DATETIME(6)  NOT NULL,
    PRIMARY KEY (id),
    INDEX idx_land_parcels_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE landowners (
    id          VARCHAR(24)  NOT NULL,
    name        VARCHAR(255) NOT NULL DEFAULT '',
    created_at  DATETIME(6)  NOT NULL,
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE agricultural_capitalists (
    id                VARCHAR(24) NOT NULL,
    capital_advanced  BIGINT      NOT NULL DEFAULT 0,
    lease_parcel_id   VARCHAR(24) NOT NULL DEFAULT '',
    created_at        DATETIME(6) NOT NULL,
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE ground_rents (
    id                     VARCHAR(24) NOT NULL,
    parcel_id              VARCHAR(24) NOT NULL DEFAULT '',
    capitalist_id          VARCHAR(24) NOT NULL DEFAULT '',
    land_owner_id          VARCHAR(24) NOT NULL DEFAULT '',
    `form`                 BIGINT      NOT NULL DEFAULT 1,
    amount_labour_minutes  BIGINT      NOT NULL DEFAULT 0,
    period_years           BIGINT      NOT NULL DEFAULT 1,
    created_at             DATETIME(6) NOT NULL,
    PRIMARY KEY (id),
    INDEX idx_ground_rents_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS ground_rents;
DROP TABLE IF EXISTS agricultural_capitalists;
DROP TABLE IF EXISTS landowners;
DROP TABLE IF EXISTS land_parcels;
