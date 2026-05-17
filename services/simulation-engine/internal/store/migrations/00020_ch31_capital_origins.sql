-- +goose Up
CREATE TABLE capital_origins (
    id                  CHAR(24)     NOT NULL PRIMARY KEY,
    historical_stage_id CHAR(24)     NOT NULL,
    source              VARCHAR(64)  NOT NULL,
    amount_pence        BIGINT       NOT NULL,
    period              VARCHAR(255) NOT NULL DEFAULT '',
    created_at          DATETIME(6)  NOT NULL,
    CONSTRAINT fk_capital_origins_stage
        FOREIGN KEY (historical_stage_id)
        REFERENCES historical_stages (id)
        ON DELETE CASCADE,
    KEY idx_capital_origins_stage (historical_stage_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE colonial_transfers (
    id                  CHAR(24)     NOT NULL PRIMARY KEY,
    historical_stage_id CHAR(24)     NOT NULL,
    `from`              VARCHAR(255) NOT NULL,
    `to`                VARCHAR(255) NOT NULL,
    value_pence         BIGINT       NOT NULL,
    method              VARCHAR(64)  NOT NULL,
    created_at          DATETIME(6)  NOT NULL,
    CONSTRAINT fk_colonial_transfers_stage
        FOREIGN KEY (historical_stage_id)
        REFERENCES historical_stages (id)
        ON DELETE CASCADE,
    KEY idx_colonial_transfers_stage (historical_stage_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE national_debts (
    id                  CHAR(24)     NOT NULL PRIMARY KEY,
    historical_stage_id CHAR(24)     NOT NULL,
    amount_pence        BIGINT       NOT NULL,
    interest_rate_bps   BIGINT       NOT NULL,
    creditor_class      VARCHAR(255) NOT NULL,
    created_at          DATETIME(6)  NOT NULL,
    CONSTRAINT fk_national_debts_stage
        FOREIGN KEY (historical_stage_id)
        REFERENCES historical_stages (id)
        ON DELETE CASCADE,
    KEY idx_national_debts_stage (historical_stage_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE protection_systems (
    id                  CHAR(24)     NOT NULL PRIMARY KEY,
    historical_stage_id CHAR(24)     NOT NULL,
    tariff_rate_bps     BIGINT       NOT NULL,
    beneficiary         VARCHAR(255) NOT NULL,
    period_start        VARCHAR(64)  NOT NULL DEFAULT '',
    period_end          VARCHAR(64)  NOT NULL DEFAULT '',
    created_at          DATETIME(6)  NOT NULL,
    CONSTRAINT fk_protection_systems_stage
        FOREIGN KEY (historical_stage_id)
        REFERENCES historical_stages (id)
        ON DELETE CASCADE,
    KEY idx_protection_systems_stage (historical_stage_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- +goose Down
DROP TABLE IF EXISTS protection_systems;
DROP TABLE IF EXISTS national_debts;
DROP TABLE IF EXISTS colonial_transfers;
DROP TABLE IF EXISTS capital_origins;
