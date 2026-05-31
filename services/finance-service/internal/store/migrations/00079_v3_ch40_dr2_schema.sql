-- +goose Up
CREATE TABLE `successive_investments` (
    `id`                VARCHAR(24) NOT NULL,
    `parcel_id`         VARCHAR(24) NOT NULL,
    `tranche_number`    INT         NOT NULL,
    `capital_bp`        BIGINT      NOT NULL,
    `output_quarters`   BIGINT      NOT NULL,
    `surplus_profit_bp` BIGINT      NOT NULL,
    `created_at`        DATETIME(6) NOT NULL,
    PRIMARY KEY (`id`),
    INDEX `idx_successive_investments_parcel_id` (`parcel_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `lease_terms` (
    `id`             VARCHAR(24) NOT NULL,
    `parcel_id`      VARCHAR(24) NOT NULL,
    `start_year`     INT         NOT NULL,
    `duration_years` INT         NOT NULL,
    `fixed_rent_bp`  BIGINT      NOT NULL,
    `created_at`     DATETIME(6) NOT NULL,
    PRIMARY KEY (`id`),
    INDEX `idx_lease_terms_parcel_id` (`parcel_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS `lease_terms`;
DROP TABLE IF EXISTS `successive_investments`;
