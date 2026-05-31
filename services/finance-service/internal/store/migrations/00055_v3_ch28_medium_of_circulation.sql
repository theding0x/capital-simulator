-- +goose Up
-- Vol. III Ch. 28 — Medium of Circulation and Capital (Tooke and Fullarton)
-- currency_observations stores CurrencyObservation persisted entities.
-- total_currency_outstanding > 0 (pence).
-- coin_function_amount + capital_transfer_amount <= total_currency_outstanding.
-- reserve_amount >= 0 (pence).
-- Column `currency_function` is intentionally absent: the split is stored as
-- two amount columns (coin_function_amount, capital_transfer_amount) rather than
-- a single function enum, to avoid the MySQL reserved word `function`.
CREATE TABLE `currency_observations` (
    `id`                          VARCHAR(24)  NOT NULL,
    `name`                        VARCHAR(255) NOT NULL DEFAULT '',
    `total_currency_outstanding`  BIGINT       NOT NULL,
    `coin_function_amount`        BIGINT       NOT NULL DEFAULT 0,
    `capital_transfer_amount`     BIGINT       NOT NULL DEFAULT 0,
    `reserve_amount`              BIGINT       NOT NULL DEFAULT 0,
    `reserve_description`         TEXT         NOT NULL,
    `description`                 TEXT         NOT NULL,
    `created_at`                  DATETIME(6)  NOT NULL,
    PRIMARY KEY (`id`),
    INDEX `idx_currency_obs_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS `currency_observations`;
