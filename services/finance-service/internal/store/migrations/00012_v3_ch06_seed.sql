-- +goose Up
-- Vol. III Ch. 6 seed — worked price fluctuations, so the price-fluctuation panel
-- comes up populated on a fresh MySQL volume. Each takes the same starting
-- capital (c = £400, v = £100, s = £100, p' = 20%) and shows how a re-pricing of
-- the raw-material element of c moves the rate of profit while s and v stay put.
--
-- Seed ID convention: 5eed00000000000000<CC><XX> where CC=06 (chapter) and XX is
-- the row sequence. These 22-char IDs never collide with the 24-char hex
-- profit.NewPriceFluctuationAnalysisID() output.
--
--   Row 0601 rise: a material-light capital, £320 fixed + £80 raw material.
--     Cotton doubles (factor 200) → c rises to £480, C' = £580; p' falls from
--     20% (2000 bp) to 17.24% (1724 bp).
--   Row 0602 fall: a material-heavy capital, £20 fixed + £380 raw material.
--     Cotton halves (factor 50) → c falls to £210, C' = £310; p' rises from
--     20% (2000 bp) to 32.26% (3225 bp, truncated).
--   Row 0603 fall: £40 fixed + £360 raw material. Cotton falls a quarter
--     (factor 75) → c falls to £310, C' = £410; p' rises from 20% (2000 bp) to
--     24.39% (2439 bp).

INSERT INTO price_fluctuation_analyses
    (id, kind, fixed_capital, original_material_capital, price_factor, variable_capital, surplus_value,
     old_profit_rate_bp, new_profit_rate_bp, created_at)
VALUES
    ('5eed000000000000000601', 'rise', 320, 80, 200, 100, 100,
     2000, 1724, '2026-05-01 00:00:00.000000'),
    ('5eed000000000000000602', 'fall', 20, 380, 50, 100, 100,
     2000, 3225, '2026-05-01 00:00:00.000000'),
    ('5eed000000000000000603', 'fall', 40, 360, 75, 100, 100,
     2000, 2439, '2026-05-01 00:00:00.000000');

-- +goose Down
DELETE FROM price_fluctuation_analyses WHERE id IN (
    '5eed000000000000000601',
    '5eed000000000000000602',
    '5eed000000000000000603'
);
