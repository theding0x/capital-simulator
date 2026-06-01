-- +goose Up
CREATE TABLE `annual_value_compositions` (
    `id` VARCHAR(24) NOT NULL,
    `constant_capital_labour_minutes` BIGINT NOT NULL,
    `variable_capital_labour_minutes` BIGINT NOT NULL,
    `surplus_value_labour_minutes` BIGINT NOT NULL,
    `total_value_labour_minutes` BIGINT NOT NULL,
    `new_value_labour_minutes` BIGINT NOT NULL,
    `created_at` DATETIME(6) NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE `revenue_limits` (
    `id` VARCHAR(24) NOT NULL,
    `composition_id` VARCHAR(24) NOT NULL,
    `max_distributable_revenue_lm` BIGINT NOT NULL,
    `constant_capital_replacement` BIGINT NOT NULL,
    `created_at` DATETIME(6) NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE `surplus_value_partitions` (
    `id` VARCHAR(24) NOT NULL,
    `composition_id` VARCHAR(24) NOT NULL,
    `total_surplus_value_lm` BIGINT NOT NULL,
    `profit_enterprise_lm` BIGINT NOT NULL,
    `interest_lm` BIGINT NOT NULL,
    `ground_rent_lm` BIGINT NOT NULL,
    `residual_surplus_lm` BIGINT NOT NULL,
    `created_at` DATETIME(6) NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS `surplus_value_partitions`;
DROP TABLE IF EXISTS `revenue_limits`;
DROP TABLE IF EXISTS `annual_value_compositions`;
