-- +goose Up
CREATE TABLE `soil_exclusion_events` (
    `id`                         VARCHAR(24) NOT NULL,
    `excluded_grade`             INT         NOT NULL,
    `new_regulating_grade`       INT         NOT NULL,
    `new_price_of_production_bp` BIGINT      NOT NULL,
    `created_at`                 DATETIME(6) NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `falling_price_dr2_outcomes` (
    `id`                    VARCHAR(24) NOT NULL,
    `falling_variant`       INT         NOT NULL,
    `exclusion_event_id`    VARCHAR(24) NOT NULL,
    `total_money_rental_bp` BIGINT      NOT NULL,
    `total_grain_rent_qtrs` BIGINT      NOT NULL,
    `total_output_quarters` BIGINT      NOT NULL,
    `total_capital_bp`      BIGINT      NOT NULL,
    `created_at`            DATETIME(6) NOT NULL,
    PRIMARY KEY (`id`),
    INDEX `idx_falling_price_dr2_outcomes_exclusion_event_id` (`exclusion_event_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `normal_capital_per_acre` (
    `id`                  VARCHAR(24)  NOT NULL,
    `period_label`        VARCHAR(100) NOT NULL,
    `capital_per_acre_bp` BIGINT       NOT NULL,
    `created_at`          DATETIME(6)  NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS `normal_capital_per_acre`;
DROP TABLE IF EXISTS `falling_price_dr2_outcomes`;
DROP TABLE IF EXISTS `soil_exclusion_events`;
