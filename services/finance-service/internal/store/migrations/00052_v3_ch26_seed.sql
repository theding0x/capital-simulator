-- +goose Up
-- Vol. III Ch. 26 seeds — Marx's money-capital accumulation examples.
-- Seed IDs start with 5eed000000000000002600..2602 (chapter 26).
--
-- Seed 2600: Post-1847 Inactivity (phase=1) — after the 1847 crisis, loanable
--            money piles up in banks; rate depressed to 1% (100 bp); cotton-mill
--            capital sits idle; gap positive (money-capital accumulates faster).
--
-- Seed 2601: Revival 1850s (phase=2) — demand for loanable capital rises;
--            supply still ample; rate moderate at 3% (300 bp); idle capital
--            re-engaged; gap narrows.
--
-- Seed 2602: Crisis Peak 1857 (phase=5) — every capitalist scrambles for means
--            of payment; loanable supply scarce; rate spikes to 10% (1000 bp);
--            no idle capital; gap negative (productive capital outpaces supply).

INSERT INTO `money_capital_accumulations`
    (`id`, `name`, `phase`, `interest_rate_bp`, `period`, `loanable_amount`, `loanable_abundant`, `idle_amount`, `idle_description`, `gap_bp`, `gap_description`, `description`, `created_at`)
VALUES
    ('5eed000000000000002600', 'Post-1847 Inactivity', 1, 100, 'post-1847 crisis', 5000000, 1, 2000000, 'idle cotton-mill capital', 500, 'money-capital accumulates faster than productive capital', 'post-crisis stagnation; loanable money piles up in banks; rate depressed', '2025-01-01 00:00:00.000000'),
    ('5eed000000000000002601', 'Revival Phase 1850s',  2, 300, '1850s revival',    4000000, 1,  500000, 'partially re-engaged capital', 100, 'gap narrows as productive capital re-activates', 'rising demand for loanable capital; supply still ample; rate moderate', '2025-01-01 00:01:00.000000'),
    ('5eed000000000000002602', 'Crisis Peak 1857',     5, 1000, '1857 crisis',      500000, 0,       0, 'no idle capital — all scrambling for means of payment', -800, 'productive capital outpaces loanable supply', 'acute demand for means of payment; rate spikes to 10%; loanable capital scarce', '2025-01-01 00:02:00.000000');

-- +goose Down
DELETE FROM `money_capital_accumulations`
WHERE `id` IN ('5eed000000000000002600', '5eed000000000000002601', '5eed000000000000002602');
