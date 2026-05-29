-- +goose Up
-- Vol. III Ch. 5 — economy analyses. Each row records one saving in constant
-- capital and the rise it produces in the rate of profit. The economy is held by
-- its kind (shared_conditions, waste_reduction, workday_extension,
-- improved_machinery, raw_material_saving) and the capital it applies to: the
-- original constant capital c, the variable capital v, the surplus-value s, and
-- the saving in c. The derived columns hold the rate of profit before the
-- economy s/(c+v), the rate after it s/(c−saving+v), and the gain between them.
-- Capital magnitudes are LabourMinutes (BIGINT); the *_bp columns are
-- dimensionless basis points (100% = 10000 bp).

CREATE TABLE IF NOT EXISTS economy_analyses (
    id                  VARCHAR(24) NOT NULL PRIMARY KEY,
    kind                VARCHAR(32) NOT NULL,
    constant_capital    BIGINT      NOT NULL,
    variable_capital    BIGINT      NOT NULL,
    surplus_value       BIGINT      NOT NULL,
    saving              BIGINT      NOT NULL,
    old_profit_rate_bp  BIGINT      NOT NULL,
    new_profit_rate_bp  BIGINT      NOT NULL,
    profit_rate_gain_bp BIGINT      NOT NULL,
    created_at          DATETIME(6) NOT NULL,
    INDEX idx_economy_analyses_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS economy_analyses;
