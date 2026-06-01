-- +goose Up
CREATE TABLE `revenue_streams` (
    `id` VARCHAR(24) NOT NULL, `source` INT NOT NULL,
    `apparent_revenue_bp` BIGINT NOT NULL, `actual_source_bp` BIGINT NOT NULL,
    `is_fetishised` TINYINT(1) NOT NULL, `created_at` DATETIME(6) NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE `trinity_formulas` (
    `id` VARCHAR(24) NOT NULL,
    `capital_stream_id` VARCHAR(24) NOT NULL, `land_stream_id` VARCHAR(24) NOT NULL, `labour_stream_id` VARCHAR(24) NOT NULL,
    `total_surplus_value_bp` BIGINT NOT NULL, `total_apparent_revenue_bp` BIGINT NOT NULL, `created_at` DATETIME(6) NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE `revenue_fetish_forms` (
    `id` VARCHAR(24) NOT NULL, `source` INT NOT NULL,
    `surface_formula` VARCHAR(255) NOT NULL, `real_relation` VARCHAR(500) NOT NULL,
    `mystification_kind` VARCHAR(255) NOT NULL, `created_at` DATETIME(6) NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS `revenue_fetish_forms`;
DROP TABLE IF EXISTS `trinity_formulas`;
DROP TABLE IF EXISTS `revenue_streams`;
