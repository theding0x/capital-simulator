-- +goose Up
-- Vol. III Ch. 3 — variation analyses. Each row records the movement of the
-- rate of profit between an initial and a changed capital decomposition under
-- Marx's equation p' = s'(v/C). Both decompositions are stored as their
-- constant component (c), variable component (v) and rate of surplus-value s'
-- (basis points, 100% = 10000 bp). The derived columns record the old and new
-- rate of profit and value-composition v/C (basis points) and whether p' moved
-- in the same proportion as s' (the §II case, v/C held constant). Capital
-- magnitudes are LabourMinutes (BIGINT); the *_bp columns are dimensionless.

CREATE TABLE IF NOT EXISTS variation_analyses (
    id                  VARCHAR(24) NOT NULL PRIMARY KEY,
    vcase               VARCHAR(32) NOT NULL,
    initial_c           BIGINT      NOT NULL,
    initial_v           BIGINT      NOT NULL,
    initial_s_rate_bp   BIGINT      NOT NULL,
    changed_c           BIGINT      NOT NULL,
    changed_v           BIGINT      NOT NULL,
    changed_s_rate_bp   BIGINT      NOT NULL,
    old_profit_rate_bp  BIGINT      NOT NULL,
    new_profit_rate_bp  BIGINT      NOT NULL,
    old_composition_bp  BIGINT      NOT NULL,
    new_composition_bp  BIGINT      NOT NULL,
    proportional_change BOOLEAN     NOT NULL,
    created_at          DATETIME(6) NOT NULL,
    INDEX idx_variation_analyses_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS variation_analyses;
