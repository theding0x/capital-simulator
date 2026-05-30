-- +goose Up
-- Vol. III Ch. 7 — the two persisted entities of the supplementary remarks.
--
-- composition_effects: two capitals carrying the same surplus-value on the same
-- variable capital but differing in their constant capital, and the gap in their
-- rates of profit. With s and v held equal the gap measures the effect of the
-- organic composition alone. The *_bp columns are basis points (100% = 10000 bp);
-- capital magnitudes are LabourMinutes (BIGINT).
--
-- magnitude_changes: a change in the absolute size of a capital and what becomes
-- of its rate of profit. factor is a percentage of the original magnitude (100 =
-- unchanged, 200 = doubled). rate_unchanged records the kernel of the refutation
-- of Rodbertus: a nominal money-value change or a proportional change at constant
-- composition leaves the rate untouched (1); only an organic change moves it (0).

CREATE TABLE IF NOT EXISTS composition_effects (
    id                 VARCHAR(24) NOT NULL PRIMARY KEY,
    s_a                BIGINT      NOT NULL,
    v_a                BIGINT      NOT NULL,
    c_a                BIGINT      NOT NULL,
    s_b                BIGINT      NOT NULL,
    v_b                BIGINT      NOT NULL,
    c_b                BIGINT      NOT NULL,
    profit_rate_a_bp   BIGINT      NOT NULL,
    profit_rate_b_bp   BIGINT      NOT NULL,
    rate_difference_bp BIGINT      NOT NULL,
    created_at         DATETIME(6) NOT NULL,
    INDEX idx_composition_effects_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS magnitude_changes (
    id               VARCHAR(24) NOT NULL PRIMARY KEY,
    kind             VARCHAR(20) NOT NULL,
    original_capital BIGINT      NOT NULL,
    original_profit  BIGINT      NOT NULL,
    factor           BIGINT      NOT NULL,
    old_rate_bp      BIGINT      NOT NULL,
    new_rate_bp      BIGINT      NOT NULL,
    rate_unchanged   TINYINT(1)  NOT NULL,
    created_at       DATETIME(6) NOT NULL,
    INDEX idx_magnitude_changes_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS magnitude_changes;
DROP TABLE IF EXISTS composition_effects;
