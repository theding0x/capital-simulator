-- +goose Up
-- Ch. 33 — The Modern Theory of Colonisation.
-- A single table records each colonial labour market. The Wakefield
-- scheme regulation outcome is stored in-row so the dashboard can
-- compare the unregulated and regulated states without joins.

CREATE TABLE colonial_labour_markets (
    id                         CHAR(24)     NOT NULL PRIMARY KEY,
    colony                     VARCHAR(255) NOT NULL,
    free_labourers             BIGINT       NOT NULL,
    annual_wage_pence          BIGINT       NOT NULL,
    land_available             TINYINT(1)   NOT NULL DEFAULT 1,
    wakefield_scheme_applied   TINYINT(1)   NOT NULL DEFAULT 0,
    independence_years         BIGINT       NOT NULL DEFAULT 0,
    surplus_labour_extractable TINYINT(1)   NOT NULL DEFAULT 0,
    created_at                 DATETIME(6)  NOT NULL,
    UNIQUE KEY uq_colonial_labour_markets_colony (colony)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- +goose Down
DROP TABLE IF EXISTS colonial_labour_markets;
