-- +goose Up
CREATE TABLE IF NOT EXISTS capital_circuits (
    id            VARCHAR(24)  NOT NULL PRIMARY KEY,
    agent_id      VARCHAR(24)  NOT NULL,
    m_advanced    BIGINT       NOT NULL,
    commodity_id  VARCHAR(255) NOT NULL,
    m_returned    BIGINT       NOT NULL DEFAULT 0,
    surplus_value BIGINT       NOT NULL DEFAULT 0,
    circuit_type  VARCHAR(20)  NOT NULL,
    created_at    DATETIME(6)  NOT NULL,
    INDEX idx_agent_id (agent_id)
);

-- +goose Down
DROP TABLE IF EXISTS capital_circuits;
