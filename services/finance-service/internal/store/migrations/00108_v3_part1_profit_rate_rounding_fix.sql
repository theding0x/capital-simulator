-- +goose Up
-- Issue #410 — Vol. III Part-I cohesion: the rate of profit p′ is now rounded
-- half-up throughout internal/profit, matching the house convention already used
-- by avgprofit and tendency (and profit's own Ch.2 rate_of_profit.go). Three
-- Part-I seed rows baked in the older truncated basis points; this brings the
-- persisted exemplars into line with what the domain now computes, so a seeded
-- analysis and a freshly POSTed one of the same capital report the same rate.
--
-- Migrations are append-only, so the original seeds (00006/00008/00012) stay
-- intact and the stored derived rates are corrected here, scoped by seed ID:
--
--   Ch.3 §I.4a variation 5eed…0301: p′ of 100c+20v at s'=100% = 10000·20/120
--     = 1666.67 → 1667 (was 1666, Marx's 16⅔%). new_composition_bp (v/C =
--     16.67%) is a different magnitude, still truncated to 1666 and left as-is —
--     #410 scopes the rate of profit, not the value-composition ratio.
--   Ch.4 Capital I turnover 5eed…0403: single-turnover p′ = 10000·500/11000
--     = 454.55 → 455 (was 454). The annual rate stays 4545 (4545.45 rounds down).
--   Ch.6 fall price-fluctuation 5eed…0602: p′ of 210c+100v = 100·10000/310
--     = 3225.81 → 3226 (was 3225).

UPDATE variation_analyses
    SET new_profit_rate_bp = 1667
    WHERE id = '5eed000000000000000301' AND new_profit_rate_bp = 1666;

UPDATE turnover_analyses
    SET single_turnover_rate_bp = 455
    WHERE id = '5eed000000000000000403' AND single_turnover_rate_bp = 454;

UPDATE price_fluctuation_analyses
    SET new_profit_rate_bp = 3226
    WHERE id = '5eed000000000000000602' AND new_profit_rate_bp = 3225;

-- +goose Down
-- Restore the truncated basis points for the three seeded rows by ID.
UPDATE variation_analyses
    SET new_profit_rate_bp = 1666
    WHERE id = '5eed000000000000000301' AND new_profit_rate_bp = 1667;

UPDATE turnover_analyses
    SET single_turnover_rate_bp = 454
    WHERE id = '5eed000000000000000403' AND single_turnover_rate_bp = 455;

UPDATE price_fluctuation_analyses
    SET new_profit_rate_bp = 3225
    WHERE id = '5eed000000000000000602' AND new_profit_rate_bp = 3226;
