-- +goose Up
-- Ch. 4 multi-agent exchange: a market-service Exchange between two owners
-- lands a complementary circuit leg on each owner via the cross-service
-- settlement webhook (#216). Loosely coupled — no FK to agents, since an
-- owner may exist only in market-service.
CREATE TABLE IF NOT EXISTS circuit_legs (
    id              VARCHAR(24) NOT NULL PRIMARY KEY,
    agent_id        VARCHAR(24) NOT NULL,
    kind            VARCHAR(16) NOT NULL,
    counterparty_id VARCHAR(24) NOT NULL,
    commodity_id    VARCHAR(255) NOT NULL,
    value_minutes   BIGINT      NOT NULL,
    exchange_id     VARCHAR(64) NOT NULL,
    created_at      DATETIME(6) NOT NULL,
    INDEX idx_circuit_legs_agent (agent_id),
    INDEX idx_circuit_legs_exchange (exchange_id)
);

-- +goose Down
DROP TABLE IF EXISTS circuit_legs;
