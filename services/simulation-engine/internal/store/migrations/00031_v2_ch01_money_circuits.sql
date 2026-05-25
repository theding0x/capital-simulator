-- +goose Up
CREATE TABLE money_circuits (
    id                   VARCHAR(24)  NOT NULL,
    agent_id             VARCHAR(24)  NOT NULL DEFAULT '',
    advance_amount       BIGINT       NOT NULL,
    advanced_at          DATETIME(6)  NOT NULL,
    moment               VARCHAR(16)  NOT NULL DEFAULT 'M',
    -- purchase phase (M—C): nullable until RecordPurchase
    labour_amount        BIGINT       NULL,
    labour_power_hours   BIGINT       NULL,
    means_amount         BIGINT       NULL,
    means_capacity_hours BIGINT       NULL,
    -- productive state (P): nullable until RecordProductive
    productive_constant  BIGINT       NULL,
    productive_variable  BIGINT       NULL,
    productive_entered_at DATETIME(6) NULL,
    -- commodity capital (C'): nullable until RecordCommodity
    commodity_id         VARCHAR(24)  NULL,
    value_original       BIGINT       NULL,
    value_surplus        BIGINT       NULL,
    -- realisation (C'--M'): nullable until RecordRealisation
    realised_pence       BIGINT       NULL,
    sold_at              DATETIME(6)  NULL,
    created_at           DATETIME(6)  NOT NULL,
    PRIMARY KEY (id),
    INDEX idx_money_circuits_agent (agent_id),
    INDEX idx_money_circuits_moment (moment)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS money_circuits;
