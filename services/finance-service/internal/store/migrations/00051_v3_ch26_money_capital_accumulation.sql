-- +goose Up
-- Vol. III Ch. 26 — Accumulation of Money-Capital and Its Influence on the Rate of Interest
-- money_capital_accumulations stores MoneyCapitalAccumulation persisted entities.
-- phase: 1=Inactivity, 2=Revival, 3=Prosperity, 4=Overproduction, 5=Crisis.
-- loanable_abundant: 1 when supply exceeds demand, 0 otherwise.
-- interest_rate_bp >= 0; loanable_amount >= 0; idle_amount >= 0.
CREATE TABLE `money_capital_accumulations` (
    `id`                VARCHAR(24)  NOT NULL,
    `name`              VARCHAR(255) NOT NULL DEFAULT '',
    `phase`             INT          NOT NULL,
    `interest_rate_bp`  BIGINT       NOT NULL DEFAULT 0,
    `period`            VARCHAR(255) NOT NULL DEFAULT '',
    `loanable_amount`   BIGINT       NOT NULL DEFAULT 0,
    `loanable_abundant` TINYINT(1)   NOT NULL DEFAULT 0,
    `idle_amount`       BIGINT       NOT NULL DEFAULT 0,
    `idle_description`  TEXT         NOT NULL,
    `gap_bp`            BIGINT       NOT NULL DEFAULT 0,
    `gap_description`   TEXT         NOT NULL,
    `description`       TEXT         NOT NULL,
    `created_at`        DATETIME(6)  NOT NULL,
    PRIMARY KEY (`id`),
    INDEX `idx_mca_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS `money_capital_accumulations`;
