-- +goose Up
-- Vol. III Ch. 25 — Credit and Fictitious Capital
-- bills_of_exchange stores BillOfExchange persisted entities.
-- FaceValue > 0; MaturityDays > 0.
CREATE TABLE `bills_of_exchange` (
    `id`            VARCHAR(24)  NOT NULL,
    `drawer`        VARCHAR(255) NOT NULL DEFAULT '',
    `drawee`        VARCHAR(255) NOT NULL DEFAULT '',
    `face_value`    BIGINT       NOT NULL,
    `maturity_days` BIGINT       NOT NULL,
    `description`   TEXT         NOT NULL,
    `created_at`    DATETIME(6)  NOT NULL,
    PRIMARY KEY (`id`),
    INDEX `idx_boe_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- fictitious_capitals stores FictitiousCapital persisted entities.
-- AnnualIncome > 0; RateBP > 0.
-- CapitalisedValue is derived: roundHalfUp(annual_income*10000, rate_bp).
-- kind: 1=bill, 2=bank-note, 3=public-debt, 4=share.
-- real_capital_exists: 0 for public-debt (consumed capital), 1 otherwise.
CREATE TABLE `fictitious_capitals` (
    `id`                   VARCHAR(24)  NOT NULL,
    `kind`                 INT          NOT NULL DEFAULT 3,
    `name`                 VARCHAR(255) NOT NULL DEFAULT '',
    `annual_income`        BIGINT       NOT NULL,
    `rate_bp`              BIGINT       NOT NULL,
    `capitalised_value`    BIGINT       NOT NULL,
    `real_capital_exists`  TINYINT(1)   NOT NULL DEFAULT 0,
    `description`          TEXT         NOT NULL,
    `created_at`           DATETIME(6)  NOT NULL,
    PRIMARY KEY (`id`),
    INDEX `idx_fc_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS `fictitious_capitals`;
DROP TABLE IF EXISTS `bills_of_exchange`;
