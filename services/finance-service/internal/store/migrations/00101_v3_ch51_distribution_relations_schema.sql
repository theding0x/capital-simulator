-- +goose Up
CREATE TABLE `production_relation_bases` (
    `id` VARCHAR(24) NOT NULL, `name` VARCHAR(100) NOT NULL,
    `characteristic_feature` VARCHAR(500) NOT NULL, `is_historically_specific` TINYINT(1) NOT NULL,
    `created_at` DATETIME(6) NOT NULL, PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE `distribution_relations` (
    `id` VARCHAR(24) NOT NULL, `form` VARCHAR(50) NOT NULL, `production_basis_id` VARCHAR(24) NOT NULL,
    `apparent_character` VARCHAR(500) NOT NULL, `real_character` VARCHAR(500) NOT NULL,
    `corresponds_to_surplus_form` VARCHAR(50) NOT NULL, `created_at` DATETIME(6) NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE `historical_transiences` (
    `id` VARCHAR(24) NOT NULL, `distribution_form_id` VARCHAR(24) NOT NULL,
    `mode_of_production` VARCHAR(100) NOT NULL, `precondition` VARCHAR(500) NOT NULL,
    `succeeding_form` VARCHAR(255) NOT NULL, `created_at` DATETIME(6) NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE `productive_force_contradictions` (
    `id` VARCHAR(24) NOT NULL, `productive_force` VARCHAR(500) NOT NULL, `social_form` VARCHAR(500) NOT NULL,
    `contradiction_nature` VARCHAR(500) NOT NULL, `crisis_signal` VARCHAR(500) NOT NULL,
    `created_at` DATETIME(6) NOT NULL, PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS `productive_force_contradictions`;
DROP TABLE IF EXISTS `historical_transiences`;
DROP TABLE IF EXISTS `distribution_relations`;
DROP TABLE IF EXISTS `production_relation_bases`;
