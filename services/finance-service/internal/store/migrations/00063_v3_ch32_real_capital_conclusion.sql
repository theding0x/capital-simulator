-- +goose Up
-- Vol. III Ch. 32 — Money-Capital and Real Capital, III (Conclusion)
-- capital_releases stores CapitalRelease persisted entities — parcels of
-- money-capital freed onto the loan market without a matching expansion of real
-- (productive) capital.
-- cause: 1 = capital released by a fall in input (raw-material) prices;
--   2 = merchants' idle money lent out between purchases and sales;
--   3 = savings of the growing rentier class placed at interest;
--   4 = revenue intended for consumption passing temporarily through money form.
-- released_amount >= 0 (pence).
-- reinvested: a release represents real accumulation only when the freed capital
--   is actually reinvested in production; otherwise it merely overstates it.
-- The embedded LoanCapitalOverstatement is flattened into overstatement_* columns:
--   overstatement_pct (basis points, >= 0 — loan capital never understates real
--   accumulation in the long run), overstatement_description.
CREATE TABLE `capital_releases` (
    `id`                        VARCHAR(24)  NOT NULL,
    `name`                      VARCHAR(255) NOT NULL DEFAULT '',
    `released_amount`           BIGINT       NOT NULL DEFAULT 0,
    `cause`                     INT          NOT NULL,
    `reinvested`                BOOLEAN      NOT NULL DEFAULT FALSE,
    `overstatement_pct`         BIGINT       NOT NULL DEFAULT 0,
    `overstatement_description` TEXT         NOT NULL,
    `description`               TEXT         NOT NULL,
    `created_at`                DATETIME(6)  NOT NULL,
    PRIMARY KEY (`id`),
    INDEX `idx_capital_releases_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS `capital_releases`;
