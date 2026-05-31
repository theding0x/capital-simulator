-- +goose Up
CREATE TABLE `agricultural_compositions` (
    `id`                      VARCHAR(24) NOT NULL,
    `constant_capital_bp`     BIGINT      NOT NULL,
    `variable_capital_bp`     BIGINT      NOT NULL,
    `composition_ratio_bp`    BIGINT      NOT NULL,
    `social_average_ratio_bp` BIGINT      NOT NULL,
    `is_below_average`        TINYINT(1)  NOT NULL,
    `created_at`              DATETIME(6) NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `value_price_gaps` (
    `id`                                 VARCHAR(24)  NOT NULL,
    `product_id`                         VARCHAR(255) NOT NULL,
    `value_labour_minutes`               BIGINT       NOT NULL,
    `price_of_production_labour_minutes` BIGINT       NOT NULL,
    `gap_labour_minutes`                 BIGINT       NOT NULL,
    `created_at`                         DATETIME(6)  NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `absolute_rents` (
    `id`                 VARCHAR(24) NOT NULL,
    `parcel_id`          VARCHAR(24) NOT NULL,
    `value_price_gap_id` VARCHAR(24) NOT NULL,
    `surplus_value_bp`   BIGINT      NOT NULL,
    `average_profit_bp`  BIGINT      NOT NULL,
    `absolute_rent_bp`   BIGINT      NOT NULL,
    `created_at`         DATETIME(6) NOT NULL,
    PRIMARY KEY (`id`),
    INDEX `idx_absolute_rents_parcel_id` (`parcel_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `absolute_rent_limits` (
    `id`             VARCHAR(24) NOT NULL,
    `max_rent_bp`    BIGINT      NOT NULL,
    `actual_rent_bp` BIGINT      NOT NULL,
    `is_at_limit`    TINYINT(1)  NOT NULL,
    `created_at`     DATETIME(6) NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS `absolute_rent_limits`;
DROP TABLE IF EXISTS `absolute_rents`;
DROP TABLE IF EXISTS `value_price_gaps`;
DROP TABLE IF EXISTS `agricultural_compositions`;
