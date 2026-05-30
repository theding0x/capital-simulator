-- +goose Up
-- Vol. III Ch. 8 — Different Compositions of Capitals in Different Branches.
--
-- production_spheres: one row per branch of production analysed. The organic
-- composition (v/C ratio) determines how much living labour the branch sets in
-- motion per £100 of total capital advanced, and therefore how much surplus-value
-- it generates and what individual rate of profit it shows. Rates are in basis
-- points (10000 = 100%); capital magnitudes are LabourMinutes (BIGINT).
--
-- id                    24-char hex (crypto/rand 96-bit)
-- name                  human-readable branch name (e.g. "Sphere A (600c+100v)")
-- c                     total capital C = c + v
-- v                     variable capital — the index of living labour
-- s_rate                rate of surplus-value s′ in bp
-- constant_capital      c = C - v (derived; stored for query convenience)
-- surplus_value         s = s′·v / 10000 (derived)
-- individual_profit_rate p′ = round_half_up(s·10000 / C) in bp (derived)
-- variable_percent       v·100 / C (derived integer %)
-- labour_power_index    (v·3600)·100 / C — living labour-minutes per £100 of C

CREATE TABLE IF NOT EXISTS production_spheres (
    id                      VARCHAR(24) NOT NULL PRIMARY KEY,
    name                    VARCHAR(120) NOT NULL,
    c                       BIGINT      NOT NULL,  -- total capital C
    v                       BIGINT      NOT NULL,  -- variable capital
    s_rate                  BIGINT      NOT NULL,  -- s' in bp
    constant_capital        BIGINT      NOT NULL,
    surplus_value           BIGINT      NOT NULL,
    individual_profit_rate  BIGINT      NOT NULL,  -- p' in bp
    variable_percent        BIGINT      NOT NULL,
    labour_power_index      BIGINT      NOT NULL,
    created_at              DATETIME(6) NOT NULL,
    INDEX idx_production_spheres_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS production_spheres;
