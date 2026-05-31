-- +goose Up
CREATE TABLE pop_surplus_profits (
    id                                VARCHAR(24)  NOT NULL,
    capital_id                        VARCHAR(255) NOT NULL DEFAULT '',
    general_price_of_production_bp    BIGINT       NOT NULL DEFAULT 0,
    individual_price_of_production_bp BIGINT       NOT NULL DEFAULT 0,
    surplus_profit_bp                 BIGINT       NOT NULL DEFAULT 0,
    created_at                        DATETIME(6)  NOT NULL,
    PRIMARY KEY (id),
    INDEX idx_pop_surplus_profits_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE monopolised_natural_forces (
    id                VARCHAR(24)  NOT NULL,
    kind              VARCHAR(255) NOT NULL DEFAULT '',
    is_monopolisable  TINYINT(1)   NOT NULL DEFAULT 1,
    parcel_id         VARCHAR(24)  NOT NULL DEFAULT '',
    created_at        DATETIME(6)  NOT NULL,
    PRIMARY KEY (id),
    INDEX idx_monopolised_natural_forces_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE differential_rents_ch38 (
    id                         VARCHAR(24) NOT NULL,
    parcel_id                  VARCHAR(24) NOT NULL DEFAULT '',
    surplus_profit_bp          BIGINT      NOT NULL DEFAULT 0,
    annual_rent_labour_minutes BIGINT      NOT NULL DEFAULT 0,
    created_at                 DATETIME(6) NOT NULL,
    PRIMARY KEY (id),
    INDEX idx_differential_rents_ch38_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE capitalised_rent_prices (
    id                               VARCHAR(24) NOT NULL,
    parcel_id                        VARCHAR(24) NOT NULL DEFAULT '',
    annual_rent_labour_minutes       BIGINT      NOT NULL DEFAULT 0,
    interest_rate_bp                 BIGINT      NOT NULL DEFAULT 0,
    capitalised_price_labour_minutes BIGINT      NOT NULL DEFAULT 0,
    created_at                       DATETIME(6) NOT NULL,
    PRIMARY KEY (id),
    INDEX idx_capitalised_rent_prices_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS capitalised_rent_prices;
DROP TABLE IF EXISTS differential_rents_ch38;
DROP TABLE IF EXISTS monopolised_natural_forces;
DROP TABLE IF EXISTS pop_surplus_profits;
