-- +goose Up
-- Vol. III Ch. 31 — Money-Capital and Real Capital, II (Continuation)
-- floating_capitals stores FloatingCapital persisted entities.
-- source: 1 = simple transformation of money into loan capital (idle post-crisis
--   industrial capital + banking economies) — INVERSELY related to real accumulation;
--   2 = transformation of capital/revenue into money that becomes loan capital —
--   ACCOMPANIES real accumulation.
-- amount >= 0 (pence).
-- The embedded LoanCapitalVelocity is flattened into velocity_* columns:
--   velocity_money_amount (pence), velocity_times_lent (>= 1),
--   velocity_loan_capital_created == velocity_money_amount × velocity_times_lent
--   (the same £20 lent 5× = £100 of loan capital).
-- The embedded CirculationReserveSaving is flattened into reserve_saving_*:
--   reserve_saving_amount (pence, >= 0), reserve_saving_description.
CREATE TABLE `floating_capitals` (
    `id`                            VARCHAR(24)  NOT NULL,
    `name`                          VARCHAR(255) NOT NULL DEFAULT '',
    `amount`                        BIGINT       NOT NULL DEFAULT 0,
    `source`                        INT          NOT NULL,
    `velocity_money_amount`         BIGINT       NOT NULL DEFAULT 0,
    `velocity_times_lent`           BIGINT       NOT NULL DEFAULT 1,
    `velocity_loan_capital_created` BIGINT       NOT NULL DEFAULT 0,
    `reserve_saving_amount`         BIGINT       NOT NULL DEFAULT 0,
    `reserve_saving_description`    TEXT         NOT NULL,
    `description`                   TEXT         NOT NULL,
    `created_at`                    DATETIME(6)  NOT NULL,
    PRIMARY KEY (`id`),
    INDEX `idx_floating_capitals_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS `floating_capitals`;
