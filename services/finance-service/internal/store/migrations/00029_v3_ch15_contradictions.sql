-- +goose Up
-- Vol. III Ch. 15 — Internal Contradictions of the Law. Completes Part III.
-- Two tables: crises and internal_contradictions. All scalar columns, no JSON.

CREATE TABLE crises (
    id                           VARCHAR(24)  NOT NULL,
    constant_capital_writedown   BIGINT       NOT NULL,  -- percent integer; 1-99
    pre_crisis_profit_rate       BIGINT       NOT NULL,  -- basis points
    post_crisis_profit_rate      BIGINT       NOT NULL,  -- basis points; derived, stored for query
    created_at                   DATETIME(6)  NOT NULL,
    PRIMARY KEY (id),
    INDEX idx_crises_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE internal_contradictions (
    id            VARCHAR(24)   NOT NULL,
    kind          VARCHAR(60)   NOT NULL,  -- ContradictionKind value
    is_coexistent TINYINT(1)    NOT NULL DEFAULT 1,
    note          VARCHAR(255)  NOT NULL DEFAULT '',
    created_at    DATETIME(6)   NOT NULL,
    PRIMARY KEY (id),
    INDEX idx_internal_contradictions_kind (kind),
    INDEX idx_internal_contradictions_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS internal_contradictions;
DROP TABLE IF EXISTS crises;
