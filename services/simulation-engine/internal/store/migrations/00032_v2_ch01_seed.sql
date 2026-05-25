-- +goose Up
-- Vol. II Ch. 1 seeds: two Spinning Mill 1871 circuits from Marx ss.II-III.
-- Seed IDs follow 5eed00000000000001<XX> (chapter 01).
--
-- SpinningMill1871 at 156% rate of surplus-value (s.II):
--   Advance M = 422 pounds (42200p): labour-leg 50 pounds / 3000h, means-leg 372 pounds / 3000h.
--   Productive state: c=372 pounds, v=50 pounds.
--   C-prime = 500 pounds yarn (value_original=42200p, value_surplus=7800p).
--   Realised M-prime = 50000p. Surplus = 7800p (78 pounds). Rate s/v = 78/50 = 156%.
INSERT INTO money_circuits
    (id, agent_id, advance_amount, advanced_at, moment,
     labour_amount, labour_power_hours, means_amount, means_capacity_hours,
     productive_constant, productive_variable, productive_entered_at,
     commodity_id, value_original, value_surplus,
     realised_pence, sold_at, created_at)
VALUES
    ('5eed000000000000010a', '', 42200, '1871-01-01 00:00:00.000000', 'M-prime',
     5000, 3000, 37200, 3000,
     37200, 5000, '1871-01-02 00:00:00.000000',
     '', 42200, 7800,
     50000, '1871-01-07 00:00:00.000000', '1871-01-01 00:00:00.000000');

-- SpinningMill1871 at 100% rate of surplus-value (s.III):
--   "If he had paid his labourers 64 pounds in wages instead of 50 pounds his
--    surplus-value would only be 64 pounds instead of 78 pounds."
--   Advance M = 43600p: labour-leg 6400p / 3000h, means-leg 37200p / 3000h.
--   Productive state: c=37200p, v=6400p.
--   C-prime = 50000p yarn (value_original=43600p, value_surplus=6400p).
--   Realised M-prime = 50000p. Surplus = 6400p. Rate s/v = 64/64 = 100%.
INSERT INTO money_circuits
    (id, agent_id, advance_amount, advanced_at, moment,
     labour_amount, labour_power_hours, means_amount, means_capacity_hours,
     productive_constant, productive_variable, productive_entered_at,
     commodity_id, value_original, value_surplus,
     realised_pence, sold_at, created_at)
VALUES
    ('5eed000000000000010b', '', 43600, '1871-01-08 00:00:00.000000', 'M-prime',
     6400, 3000, 37200, 3000,
     37200, 6400, '1871-01-09 00:00:00.000000',
     '', 43600, 6400,
     50000, '1871-01-14 00:00:00.000000', '1871-01-08 00:00:00.000000');

-- +goose Down
DELETE FROM money_circuits WHERE id IN ('5eed000000000000010a', '5eed000000000000010b');
