-- +goose Up
CREATE TABLE historical_stages (
    id          CHAR(24)     NOT NULL PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    description TEXT         NOT NULL,
    created_at  DATETIME(6)  NOT NULL,
    UNIQUE KEY uniq_historical_stages_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE primitive_accumulations (
    id                     CHAR(24)     NOT NULL PRIMARY KEY,
    stage_id               CHAR(24)     NOT NULL,
    ordinal                INT          NOT NULL,
    period                 VARCHAR(255) NOT NULL,
    method                 VARCHAR(255) NOT NULL,
    labourers_expropriated BIGINT       NOT NULL,
    capital_formed         BIGINT       NOT NULL,
    CONSTRAINT fk_primitive_accumulations_stage
        FOREIGN KEY (stage_id) REFERENCES historical_stages(id)
        ON DELETE CASCADE,
    KEY idx_primitive_accumulations_stage (stage_id, ordinal)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- +goose Down
DROP TABLE IF EXISTS primitive_accumulations;
DROP TABLE IF EXISTS historical_stages;
