-- +goose Up
-- Vol. III Ch. 21 — Interest-Bearing Capital (opens Part V)
-- interest_bearing_capitals stores InterestBearingCapital persisted entities.
-- Interest is a portion of average profit: interest_earned = roundHalfUp(money_advanced * interest_rate_bp, 10000).
-- new_value_created is always 0 — interest-bearing capital creates no new value.
CREATE TABLE `interest_bearing_capitals` (
    `id`                VARCHAR(24)   NOT NULL,
    `money_advanced`    BIGINT        NOT NULL,
    `interest_rate_bp`  BIGINT        NOT NULL,
    `loan_term_days`    BIGINT        NOT NULL,
    `interest_earned`   BIGINT        NOT NULL,
    `money_returned`    BIGINT        NOT NULL,
    `new_value_created` BIGINT        NOT NULL DEFAULT 0,
    `created_at`        DATETIME(6)   NOT NULL,
    PRIMARY KEY (`id`),
    INDEX `idx_ibc_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS `interest_bearing_capitals`;
