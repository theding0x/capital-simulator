-- +goose Up
-- Seed day_wages for the four countries seeded in national_intensities but missing
-- from day_wages in 00021_ch22_seed.sql. Wages drawn from Marx's Ch.22 comparative
-- data (approximate values in pence at 1866 exchange rates).

INSERT INTO day_wages (country_code, nominal_pence, working_day_minutes) VALUES
    ('RU', 12, 780),   -- Russia: low nominal, long working day, intensity 0.8
    ('AT', 18, 720),   -- Austria: moderate nominal, standard day, intensity 0.9
    ('BE', 26, 720),   -- Belgium: moderate-high nominal, standard day, intensity 1.1
    ('CH', 30, 600);   -- Switzerland: relatively high nominal, shorter day, intensity 1.1

-- +goose Down
DELETE FROM day_wages
WHERE country_code IN ('RU', 'AT', 'BE', 'CH');
