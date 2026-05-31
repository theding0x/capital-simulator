-- +goose Up
-- Vol. III Ch. 29 — Component Parts of Bank Capital
-- bank_capitals stores BankCapital persisted entities.
--   TotalCapital == CashAmount + SecuritiesAmount enforced by application.
--   components_json is a JSON-encoded array of BankCapitalComponent
--   ({amount, kind, description}); kind is 1=cash,2=bill,3=bond,4=stock,5=mortgage.
CREATE TABLE `bank_capitals` (
    `id`                VARCHAR(24)  NOT NULL,
    `name`              VARCHAR(255) NOT NULL DEFAULT '',
    `cash_amount`       BIGINT       NOT NULL DEFAULT 0,
    `securities_amount` BIGINT       NOT NULL DEFAULT 0,
    `total_capital`     BIGINT       NOT NULL DEFAULT 0,
    `components_json`   MEDIUMTEXT   NOT NULL,
    `description`       TEXT         NOT NULL,
    `created_at`        DATETIME(6)  NOT NULL,
    PRIMARY KEY (`id`),
    INDEX `idx_bankcap_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- fictitious_capital_valuations stores FictitiousCapitalValuation persisted entities.
--   InterestRateBP > 0.
--   MarketValue is derived: roundHalfUp(annual_income*10000, interest_rate_bp);
--   it moves inversely with the rate of interest.
--   dividend_rate_bp is optional (for stocks) and does not alter market_value.
CREATE TABLE `fictitious_capital_valuations` (
    `id`               VARCHAR(24)  NOT NULL,
    `name`             VARCHAR(255) NOT NULL DEFAULT '',
    `nominal_value`    BIGINT       NOT NULL DEFAULT 0,
    `annual_income`    BIGINT       NOT NULL DEFAULT 0,
    `interest_rate_bp` BIGINT       NOT NULL,
    `dividend_rate_bp` BIGINT       NOT NULL DEFAULT 0,
    `market_value`     BIGINT       NOT NULL DEFAULT 0,
    `description`      TEXT         NOT NULL,
    `created_at`       DATETIME(6)  NOT NULL,
    PRIMARY KEY (`id`),
    INDEX `idx_fcv_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS `fictitious_capital_valuations`;
DROP TABLE IF EXISTS `bank_capitals`;
