-- +goose Up
CREATE TABLE `building_site_rents` (
    `id` VARCHAR(24) NOT NULL, `parcel_id` VARCHAR(24) NOT NULL, `location_grade` INT NOT NULL,
    `annual_rent_labour_minutes` BIGINT NOT NULL, `lease_term_years` INT NOT NULL, `created_at` DATETIME(6) NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE `mining_rents` (
    `id` VARCHAR(24) NOT NULL, `parcel_id` VARCHAR(24) NOT NULL, `ore_grade` INT NOT NULL,
    `annual_rent_labour_minutes` BIGINT NOT NULL, `surplus_profit_bp` BIGINT NOT NULL, `created_at` DATETIME(6) NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE `monopoly_price_rents` (
    `id` VARCHAR(24) NOT NULL, `parcel_id` VARCHAR(24) NOT NULL, `product_kind` VARCHAR(255) NOT NULL,
    `value_labour_minutes` BIGINT NOT NULL, `monopoly_sale_price_labour_minutes` BIGINT NOT NULL,
    `monopoly_rent_labour_minutes` BIGINT NOT NULL, `created_at` DATETIME(6) NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE `land_price_scenarios` (
    `id` VARCHAR(24) NOT NULL, `kind` INT NOT NULL, `parcel_id` VARCHAR(24) NOT NULL,
    `old_annual_rent_bp` BIGINT NOT NULL, `new_annual_rent_bp` BIGINT NOT NULL,
    `old_interest_rate_bp` BIGINT NOT NULL, `new_interest_rate_bp` BIGINT NOT NULL,
    `old_land_price_bp` BIGINT NOT NULL, `new_land_price_bp` BIGINT NOT NULL, `created_at` DATETIME(6) NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS `land_price_scenarios`;
DROP TABLE IF EXISTS `monopoly_price_rents`;
DROP TABLE IF EXISTS `mining_rents`;
DROP TABLE IF EXISTS `building_site_rents`;
