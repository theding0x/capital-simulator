-- +goose Up
CREATE TABLE `rent_historical_forms` (
    `id` VARCHAR(24) NOT NULL, `stage` INT NOT NULL, `name` VARCHAR(100) NOT NULL,
    `coercion_kind` VARCHAR(100) NOT NULL, `surplus_form` VARCHAR(100) NOT NULL,
    `example_region` VARCHAR(255) NOT NULL, `example_period` VARCHAR(100) NOT NULL, `created_at` DATETIME(6) NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE `historical_rent_transitions` (
    `id` VARCHAR(24) NOT NULL, `from_stage` INT NOT NULL, `to_stage` INT NOT NULL,
    `preconditions` VARCHAR(500) NOT NULL, `driving_force` VARCHAR(255) NOT NULL, `created_at` DATETIME(6) NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE `small_peasant_productions` (
    `id` VARCHAR(24) NOT NULL, `parcel_id` VARCHAR(24) NOT NULL, `total_income_bp` BIGINT NOT NULL,
    `rent_component_bp` BIGINT NOT NULL, `profit_component_bp` BIGINT NOT NULL, `wages_component_bp` BIGINT NOT NULL,
    `is_conflated` TINYINT(1) NOT NULL, `created_at` DATETIME(6) NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS `small_peasant_productions`;
DROP TABLE IF EXISTS `historical_rent_transitions`;
DROP TABLE IF EXISTS `rent_historical_forms`;
