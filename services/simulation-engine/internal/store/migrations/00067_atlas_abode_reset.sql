-- +goose Up
-- Atlas fix: the General-Law step previously destroyed variable capital each
-- period (machinery displacing v→c), driving v→0 — a degenerate collapse where
-- Σv=0, the wage pressure pegged at 100%, and the reserve army swallowed the
-- whole workforce. The step is now fixed (variable capital grows via accumulation;
-- the reserve army grows with the rising organic composition, never by destroying
-- v). A singleton stranded near v=0 cannot self-recover (no surplus → no
-- accumulation), so reset it to fresh seed values and clear the immiseration
-- series so the chart restarts clean. The seed marginal composition is set above
-- the 2:1 stock so the composition rises from the very first period.
UPDATE abode_state SET
    period                  = 0,
    constant_pence          = 600000,
    variable_pence          = 300000,
    base_wage_pence         = 2500,
    worker_supply           = 150,
    surplus_rate_base_bp    = 10000,
    accumulation_rate_bp    = 5000,
    marginal_composition_bp = 8000,
    displacement_rate_bp    = 1800,
    productivity_growth_bp  = 500,
    population_growth_bp    = 150
WHERE id = '5eed000000000000abode1';

DELETE FROM general_law_periods;

-- +goose Down
-- One-way data correction (heals a degenerate persisted state); no meaningful rollback.
SELECT 1;
