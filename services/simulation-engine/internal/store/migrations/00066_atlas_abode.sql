-- +goose Up
-- Atlas Slice 2: the hidden abode. abode_state is the single evolving aggregate
-- class relation (Vol. I Ch. 25); general_law_periods is the immiseration
-- time-series the General-Law ticker appends to each scheduler pass.
CREATE TABLE abode_state (
    id                      VARCHAR(64) NOT NULL,
    period                  BIGINT      NOT NULL,
    constant_pence          BIGINT      NOT NULL,
    variable_pence          BIGINT      NOT NULL,
    base_wage_pence         BIGINT      NOT NULL,
    worker_supply           BIGINT      NOT NULL,
    surplus_rate_base_bp    BIGINT      NOT NULL,
    accumulation_rate_bp    BIGINT      NOT NULL,
    marginal_composition_bp BIGINT      NOT NULL,
    displacement_rate_bp    BIGINT      NOT NULL,
    productivity_growth_bp  BIGINT      NOT NULL,
    population_growth_bp    BIGINT      NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE general_law_periods (
    id                      VARCHAR(64) NOT NULL,
    period                  BIGINT      NOT NULL,
    wage_pence              BIGINT      NOT NULL,
    rate_of_exploitation_bp BIGINT      NOT NULL,
    reserve_army_count      BIGINT      NOT NULL,
    organic_composition_bp  BIGINT      NOT NULL,
    employed_count          BIGINT      NOT NULL,
    created_at              DATETIME(6) NOT NULL,
    PRIMARY KEY (id),
    KEY idx_glp_period (period)
);

-- The singleton abode_state, mirroring simulation.NewAbodeState().
INSERT INTO abode_state
    (id, period, constant_pence, variable_pence, base_wage_pence, worker_supply,
     surplus_rate_base_bp, accumulation_rate_bp, marginal_composition_bp,
     displacement_rate_bp, productivity_growth_bp, population_growth_bp)
VALUES
    ('5eed000000000000abode1', 0, 600000, 300000, 2500, 150,
     10000, 5000, 6667, 1800, 500, 150);

-- +goose Down
DROP TABLE IF EXISTS general_law_periods;
DROP TABLE IF EXISTS abode_state;
