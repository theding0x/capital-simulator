-- +goose Up
-- Vol. I Ch. 3 — money-form records: hoards, payment obligations, world-money transfers.

CREATE TABLE IF NOT EXISTS hoards (
    id          VARCHAR(24) NOT NULL PRIMARY KEY,
    owner_id    VARCHAR(24) NOT NULL,
    amount      BIGINT      NOT NULL,
    created_at  DATETIME(6) NOT NULL,
    INDEX idx_hoards_owner (owner_id),
    INDEX idx_hoards_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS payment_obligations (
    id          VARCHAR(24) NOT NULL PRIMARY KEY,
    creditor_id VARCHAR(24) NOT NULL,
    debtor_id   VARCHAR(24) NOT NULL,
    amount      BIGINT      NOT NULL,
    created_at  DATETIME(6) NOT NULL,
    due_at      DATETIME(6) NOT NULL,
    paid_at     DATETIME(6) NULL,
    INDEX idx_payment_obligations_creditor (creditor_id),
    INDEX idx_payment_obligations_debtor (debtor_id),
    INDEX idx_payment_obligations_due (due_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS world_money_transfers (
    id           VARCHAR(24) NOT NULL PRIMARY KEY,
    sender_id    VARCHAR(24) NOT NULL,
    receiver_id  VARCHAR(24) NOT NULL,
    gold_mg      BIGINT      NOT NULL,
    created_at   DATETIME(6) NOT NULL,
    INDEX idx_world_money_sender (sender_id),
    INDEX idx_world_money_receiver (receiver_id),
    INDEX idx_world_money_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS world_money_transfers;
DROP TABLE IF EXISTS payment_obligations;
DROP TABLE IF EXISTS hoards;
