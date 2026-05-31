-- +goose Up
-- Vol. III Ch. 30 — Money-Capital and Real Capital, I
-- real_capital_accumulations stores RealCapitalAccumulation persisted entities.
-- phase: 1=Inactivity, 2=Revival, 3=Prosperity, 4=Overproduction, 5=Crisis.
-- money_capital_growth_bp / real_capital_growth_bp may each be positive OR
-- negative; there is no equality invariant between them (the whole point of the
-- chapter is that the two accumulations can diverge).
-- corresponds: 1 when money- and real-capital accumulation align, 0 otherwise.
-- reserve_capital >= 0 (pence) — the elastic limit of the credit system.
-- commercial_credit_contracts: 1 when commercial credit is contracting.
-- cause_is_money_scarcity: 1 if the contraction is blamed on a money shortage;
--   Marx insists the crisis cause is usually blocked returns + lost confidence.
CREATE TABLE `real_capital_accumulations` (
    `id`                          VARCHAR(24)  NOT NULL,
    `name`                        VARCHAR(255) NOT NULL DEFAULT '',
    `phase`                       INT          NOT NULL,
    `money_capital_growth_bp`     BIGINT       NOT NULL DEFAULT 0,
    `real_capital_growth_bp`      BIGINT       NOT NULL DEFAULT 0,
    `corresponds`                 TINYINT(1)   NOT NULL DEFAULT 0,
    `correspondence_explanation`  TEXT         NOT NULL,
    `reserve_capital`             BIGINT       NOT NULL DEFAULT 0,
    `credit_limit_description`    TEXT         NOT NULL,
    `commercial_credit_contracts` TINYINT(1)   NOT NULL DEFAULT 0,
    `cause_is_money_scarcity`     TINYINT(1)   NOT NULL DEFAULT 0,
    `description`                 TEXT         NOT NULL,
    `created_at`                  DATETIME(6)  NOT NULL,
    PRIMARY KEY (`id`),
    INDEX `idx_rca_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS `real_capital_accumulations`;
