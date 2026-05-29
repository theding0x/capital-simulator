-- +goose Up
CREATE TABLE wage_statutes (
    id                   CHAR(24)     NOT NULL PRIMARY KEY,
    historical_stage_id  CHAR(24)     NOT NULL,
    period               VARCHAR(255) NOT NULL,
    jurisdiction         VARCHAR(255) NOT NULL,
    max_wage_pence       BIGINT       NOT NULL,
    min_wage_pence       BIGINT       NOT NULL,
    enforcement_penalty  VARCHAR(255) NOT NULL,
    created_at           DATETIME(6)  NOT NULL,
    KEY idx_wage_statutes_stage (historical_stage_id),
    KEY idx_wage_statutes_created (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE vagrancy_laws (
    id                   CHAR(24)     NOT NULL PRIMARY KEY,
    historical_stage_id  CHAR(24)     NOT NULL,
    period               VARCHAR(255) NOT NULL,
    jurisdiction         VARCHAR(255) NOT NULL,
    punishment           VARCHAR(255) NOT NULL,
    target_population    VARCHAR(255) NOT NULL,
    created_at           DATETIME(6)  NOT NULL,
    KEY idx_vagrancy_laws_stage (historical_stage_id),
    KEY idx_vagrancy_laws_created (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- +goose Down
DROP TABLE IF EXISTS vagrancy_laws;
DROP TABLE IF EXISTS wage_statutes;