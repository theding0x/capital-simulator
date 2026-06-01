-- +goose Up
CREATE TABLE `competition_illusions` (
    `id` VARCHAR(24) NOT NULL, `name` VARCHAR(255) NOT NULL,
    `surface_appearance` VARCHAR(500) NOT NULL, `real_relation` VARCHAR(500) NOT NULL,
    `produced_by` VARCHAR(100) NOT NULL, `created_at` DATETIME(6) NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE `cost_price_fetishes` (
    `id` VARCHAR(24) NOT NULL, `constant_capital_bp` BIGINT NOT NULL, `variable_capital_bp` BIGINT NOT NULL,
    `cost_price_bp` BIGINT NOT NULL, `profit_bp` BIGINT NOT NULL, `actual_value_bp` BIGINT NOT NULL,
    `created_at` DATETIME(6) NOT NULL, PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE `annual_revenue_accounts` (
    `id` VARCHAR(24) NOT NULL, `new_value_created_lm` BIGINT NOT NULL, `wages_share_lm` BIGINT NOT NULL,
    `profit_and_rent_share_lm` BIGINT NOT NULL, `constant_capital_lm` BIGINT NOT NULL, `total_social_product_lm` BIGINT NOT NULL,
    `created_at` DATETIME(6) NOT NULL, PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE `value_appearance_contrasts` (
    `id` VARCHAR(24) NOT NULL, `illusion_id` VARCHAR(24) NOT NULL,
    `surface` VARCHAR(500) NOT NULL, `reality` VARCHAR(500) NOT NULL,
    `created_at` DATETIME(6) NOT NULL, PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS `value_appearance_contrasts`;
DROP TABLE IF EXISTS `annual_revenue_accounts`;
DROP TABLE IF EXISTS `cost_price_fetishes`;
DROP TABLE IF EXISTS `competition_illusions`;
