-- +goose Up
CREATE TABLE `rent_on_worst_soils` (
    `id`                         VARCHAR(24) NOT NULL,
    `parcel_id`                  VARCHAR(24) NOT NULL,
    `mechanism`                  INT         NOT NULL,
    `old_price_of_production_bp` BIGINT      NOT NULL,
    `new_price_of_production_bp` BIGINT      NOT NULL,
    `rent_per_acre_bp`           BIGINT      NOT NULL,
    `created_at`                 DATETIME(6) NOT NULL,
    PRIMARY KEY (`id`),
    INDEX `idx_rent_on_worst_soils_parcel_id` (`parcel_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `soil_improvement_capitals` (
    `id`                   VARCHAR(24) NOT NULL,
    `parcel_id`            VARCHAR(24) NOT NULL,
    `capital_invested_bp`  BIGINT      NOT NULL,
    `productivity_gain_bp` BIGINT      NOT NULL,
    `amortisation_years`   INT         NOT NULL,
    `amortised_years`      INT         NOT NULL,
    `created_at`           DATETIME(6) NOT NULL,
    PRIMARY KEY (`id`),
    INDEX `idx_soil_improvement_capitals_parcel_id` (`parcel_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `regulating_price_shifts` (
    `id`           VARCHAR(24) NOT NULL,
    `from_grade`   INT         NOT NULL,
    `to_grade`     INT         NOT NULL,
    `old_price_bp` BIGINT      NOT NULL,
    `new_price_bp` BIGINT      NOT NULL,
    `created_at`   DATETIME(6) NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS `regulating_price_shifts`;
DROP TABLE IF EXISTS `soil_improvement_capitals`;
DROP TABLE IF EXISTS `rent_on_worst_soils`;
