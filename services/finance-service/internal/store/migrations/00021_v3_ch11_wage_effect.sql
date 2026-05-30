-- +goose Up
-- Vol. III Ch. 11 — Effects of General Wage Fluctuations on Prices of Production.
-- One table: wage_effect_analyses. Prices stored in centi-units (×100).

CREATE TABLE wage_effect_analyses (
    id                   VARCHAR(24)  NOT NULL,
    base_constant        BIGINT       NOT NULL, -- c of the average capital, whole units
    base_variable        BIGINT       NOT NULL, -- v of the average capital, whole units
    s_rate               BIGINT       NOT NULL, -- s' in basis points (10000 = 100%)
    wage_factor          BIGINT       NOT NULL, -- percent: 100 = no change, 125 = +25%
    old_general_rate     BIGINT       NOT NULL, -- p' before wage change, basis points
    new_general_rate     BIGINT       NOT NULL, -- p' after wage change, basis points
    kind                 VARCHAR(8)   NOT NULL, -- 'rise' | 'fall' | ''
    average_outcome_json JSON         NOT NULL, -- SphereWageOutcome for the base sphere
    outcomes_json        JSON         NOT NULL, -- []SphereWageOutcome for all spheres
    created_at           DATETIME(6)  NOT NULL,
    PRIMARY KEY (id),
    INDEX idx_wage_effect_analyses_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS wage_effect_analyses;
