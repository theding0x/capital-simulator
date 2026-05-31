-- +goose Up
CREATE TABLE `dr2_outcomes` (
    `id`                        VARCHAR(24) NOT NULL,
    `parcel_id`                 VARCHAR(24) NOT NULL,
    `productivity_case`         INT         NOT NULL,
    `additional_capital_bp`     BIGINT      NOT NULL,
    `new_rent_per_acre_bp`      BIGINT      NOT NULL,
    `total_rental_bp`           BIGINT      NOT NULL,
    `rate_of_surplus_profit_bp` BIGINT      NOT NULL,
    `created_at`                DATETIME(6) NOT NULL,
    PRIMARY KEY (`id`),
    INDEX `idx_dr2_outcomes_parcel_id` (`parcel_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `intensive_extensive_comparisons` (
    `id`                         VARCHAR(24) NOT NULL,
    `intensive_rent_per_acre_bp` BIGINT      NOT NULL,
    `extensive_rent_per_acre_bp` BIGINT      NOT NULL,
    `total_capital_bp`           BIGINT      NOT NULL,
    `total_rental_same_bp`       BIGINT      NOT NULL,
    `created_at`                 DATETIME(6) NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS `intensive_extensive_comparisons`;
DROP TABLE IF EXISTS `dr2_outcomes`;
