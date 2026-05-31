-- +goose Up
-- Vol. III Ch. 24 — Externalisation of the Relations of Capital in the Form of Interest-Bearing Capital
-- compound_interest_schedules stores CompoundInterestSchedule persisted entities.
-- Principal > 0; RateBP > 0; PeriodYears > 0.
-- FinalValue is derived by integer compound iteration: v_{k+1} = roundHalfUp(v_k*(10000+rate_bp), 10000).
-- NewValueCreated is always 0 (fetish form invariant).
-- fetish_form: 1=SimpleInterest, 2=CompoundInterest, 3=MM.
CREATE TABLE `compound_interest_schedules` (
    `id`               VARCHAR(24)  NOT NULL,
    `principal`        BIGINT       NOT NULL,
    `rate_bp`          BIGINT       NOT NULL,
    `period_years`     BIGINT       NOT NULL,
    `final_value`      BIGINT       NOT NULL,
    `new_value_created` BIGINT      NOT NULL DEFAULT 0,
    `fetish_form`      INT          NOT NULL DEFAULT 2,
    `created_at`       DATETIME(6)  NOT NULL,
    PRIMARY KEY (`id`),
    INDEX `idx_cis_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS `compound_interest_schedules`;
