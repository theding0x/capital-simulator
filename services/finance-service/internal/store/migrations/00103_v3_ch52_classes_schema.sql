-- +goose Up
CREATE TABLE `social_classes` (
    `id` VARCHAR(24) NOT NULL, `name` VARCHAR(100) NOT NULL,
    `relation_to_means_of_production` VARCHAR(500) NOT NULL, `primary_revenue_form` VARCHAR(50) NOT NULL,
    `historical_precondition` VARCHAR(500) NOT NULL, `created_at` DATETIME(6) NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE `class_income_sources` (
    `id` VARCHAR(24) NOT NULL, `class_id` VARCHAR(24) NOT NULL, `revenue_form` VARCHAR(50) NOT NULL,
    `surplus_value_component_bp` BIGINT NOT NULL, `created_at` DATETIME(6) NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE `class_tendencies` (
    `id` VARCHAR(24) NOT NULL, `class_id` VARCHAR(24) NOT NULL, `description` VARCHAR(500) NOT NULL,
    `direction` VARCHAR(50) NOT NULL, `created_at` DATETIME(6) NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE `class_ambiguities` (
    `id` VARCHAR(24) NOT NULL, `intermediate_group` VARCHAR(255) NOT NULL, `income_source` VARCHAR(255) NOT NULL,
    `why_not_a_class` VARCHAR(255) NOT NULL, `created_at` DATETIME(6) NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS `class_ambiguities`;
DROP TABLE IF EXISTS `class_tendencies`;
DROP TABLE IF EXISTS `class_income_sources`;
DROP TABLE IF EXISTS `social_classes`;
