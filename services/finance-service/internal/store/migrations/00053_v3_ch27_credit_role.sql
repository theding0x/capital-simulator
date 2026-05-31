-- +goose Up
-- Vol. III Ch. 27 — The Role of Credit in Capitalist Production
-- stock_companies: joint-stock companies (abolition of private capital within capitalism).
--   TotalCapital == FixedCapital + FloatingCapital enforced by application.
--   manager_wage_minutes > 0: wages of superintendence, NOT profit of enterprise.
-- cooperative_factories: worker-owned transitional form; worker_owned always 1.
CREATE TABLE `stock_companies` (
    `id`                   VARCHAR(24)  NOT NULL,
    `name`                 VARCHAR(255) NOT NULL DEFAULT '',
    `fixed_capital`        BIGINT       NOT NULL DEFAULT 0,
    `floating_capital`     BIGINT       NOT NULL DEFAULT 0,
    `total_capital`        BIGINT       NOT NULL DEFAULT 0,
    `manager_name`         VARCHAR(255) NOT NULL DEFAULT '',
    `manager_title`        VARCHAR(255) NOT NULL DEFAULT '',
    `manager_wage_minutes` BIGINT       NOT NULL DEFAULT 0,
    `description`          TEXT         NOT NULL,
    `created_at`           DATETIME(6)  NOT NULL,
    PRIMARY KEY (`id`),
    INDEX `idx_sc_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `cooperative_factories` (
    `id`               VARCHAR(24)  NOT NULL,
    `name`             VARCHAR(255) NOT NULL DEFAULT '',
    `worker_owned`     TINYINT(1)   NOT NULL DEFAULT 1,
    `worker_count`     BIGINT       NOT NULL DEFAULT 0,
    `capital_advanced` BIGINT       NOT NULL DEFAULT 0,
    `description`      TEXT         NOT NULL,
    `created_at`       DATETIME(6)  NOT NULL,
    PRIMARY KEY (`id`),
    INDEX `idx_cf_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS `cooperative_factories`;
DROP TABLE IF EXISTS `stock_companies`;
