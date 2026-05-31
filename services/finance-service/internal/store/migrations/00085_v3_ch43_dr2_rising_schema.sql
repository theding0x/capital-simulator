-- +goose Up
CREATE TABLE `rising_price_dr2_outcomes` (
    `id`                         VARCHAR(24) NOT NULL,
    `rising_variant`             INT         NOT NULL,
    `total_capital_shillings`    BIGINT      NOT NULL,
    `total_rent_shillings`       BIGINT      NOT NULL,
    `total_output_bushels`       BIGINT      NOT NULL,
    `price_per_bushel_shillings` BIGINT      NOT NULL,
    `created_at`                 DATETIME(6) NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `rent_sequences` (
    `id`             VARCHAR(24) NOT NULL,
    `base_rent`      BIGINT      NOT NULL,
    `rent_increment` BIGINT      NOT NULL,
    `num_grades`     INT         NOT NULL,
    `created_at`     DATETIME(6) NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `engels_table_entries` (
    `id`                VARCHAR(24) NOT NULL,
    `table_number`      INT         NOT NULL,
    `soil_grade_label`  VARCHAR(10) NOT NULL,
    `output_bushels`    BIGINT      NOT NULL,
    `rent_shillings`    BIGINT      NOT NULL,
    `capital_shillings` BIGINT      NOT NULL,
    `created_at`        DATETIME(6) NOT NULL,
    PRIMARY KEY (`id`),
    INDEX `idx_engels_table_entries_table_number` (`table_number`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS `engels_table_entries`;
DROP TABLE IF EXISTS `rent_sequences`;
DROP TABLE IF EXISTS `rising_price_dr2_outcomes`;
