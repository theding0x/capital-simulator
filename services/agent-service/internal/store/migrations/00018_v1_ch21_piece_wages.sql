-- +goose Up
CREATE TABLE piece_wages (
    id            CHAR(24)    NOT NULL PRIMARY KEY,
    agent_id      CHAR(24)    NOT NULL,
    price_pence   BIGINT      NOT NULL,
    normal_output BIGINT      NOT NULL,
    created_at    DATETIME(6) NOT NULL,
    UNIQUE KEY uk_agent_piece_wage (agent_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE sub_contracts (
    id                   CHAR(24)    NOT NULL PRIMARY KEY,
    head_labourer_id     CHAR(24)    NOT NULL,
    assistant_ids        JSON        NOT NULL,
    piece_rate_pence     BIGINT      NOT NULL,
    assistant_rate_pence BIGINT      NOT NULL,
    created_at           DATETIME(6) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS sub_contracts;
DROP TABLE IF EXISTS piece_wages;
