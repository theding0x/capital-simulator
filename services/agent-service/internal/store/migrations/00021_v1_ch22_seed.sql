-- +goose Up
-- Seed data for Capital Vol. I, Ch. 22 — National Differences of Wages.
--
-- Intensities: GB = 2× international average (Cowell's "England virtually lower to capitalist");
-- FR ≈ 1.2 (moderate); DE (Prussia) = 1.0 (international baseline).
--
-- Day wages (in pence): GB 48d / 600 min; FR 36d / 720 min; DE 24d / 720 min.
-- Spec fixture: StandardiseWage(DayWage{FR,36,720}, 600) == StandardisedWage{FR,30}
--
-- Spindle ratios from Redgrave's 1866 factory inspectorate data.

INSERT INTO national_intensities (country_code, factor) VALUES
    ('GB', 2.0),
    ('FR', 1.2),
    ('DE', 1.0),
    ('RU', 0.8),
    ('AT', 0.9),
    ('BE', 1.1),
    ('CH', 1.1);

INSERT INTO day_wages (country_code, nominal_pence, working_day_minutes) VALUES
    ('GB', 48, 600),
    ('FR', 36, 720),
    ('DE', 24, 720);

INSERT INTO spindle_ratios (country_code, spindles_per_worker) VALUES
    ('GB', 74),
    ('FR', 14),
    ('RU', 28),
    ('DE', 37),
    ('AT', 49),
    ('BE', 50),
    ('CH', 55);

-- +goose Down
DELETE FROM spindle_ratios
WHERE country_code IN ('GB', 'FR', 'RU', 'DE', 'AT', 'BE', 'CH');

DELETE FROM day_wages
WHERE country_code IN ('GB', 'FR', 'DE');

DELETE FROM national_intensities
WHERE country_code IN ('GB', 'FR', 'DE', 'RU', 'AT', 'BE', 'CH');
