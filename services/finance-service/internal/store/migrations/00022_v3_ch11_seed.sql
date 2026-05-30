-- +goose Up
-- Vol. III Ch. 11 seed — Marx's §main wage-rise analysis (base 80c+20v, s'=100%,
-- factor=125, three spheres: average/lower/higher) so the Ch. 11 panel comes up
-- populated on a fresh MySQL volume.
--
-- Seed ID convention: 5eed00000000000000<CC><XX> where CC=11, XX=row.
-- These 22-char IDs never collide with the 24-char hex output of New*ID().
--
-- All price/cost fields in centi-units (×100):
--   average 80c+20v: old_price=12000 new_price=12000 (unchanged)
--   lower   50c+50v: old_price=12000 new_price=12857 (rises)
--   higher  92c+8v:  old_price=12000 new_price=11657 (falls)

INSERT INTO wage_effect_analyses
    (id, base_constant, base_variable, s_rate, wage_factor,
     old_general_rate, new_general_rate, kind,
     average_outcome_json, outcomes_json, created_at)
VALUES (
    '5eed000000000000001101',
    80, 20, 10000, 125,
    2000, 1429, 'rise',
    '{"sphere_name":"average","constant_capital":80,"variable_capital":20,"wage_factor":125,"composition":"average","old_price":12000,"new_price":12000,"price_change":0,"direction":"unchanged"}',
    '[{"sphere_name":"average","constant_capital":80,"variable_capital":20,"wage_factor":125,"composition":"average","old_price":12000,"new_price":12000,"price_change":0,"direction":"unchanged"},{"sphere_name":"lower","constant_capital":50,"variable_capital":50,"wage_factor":125,"composition":"lower","old_price":12000,"new_price":12857,"price_change":857,"direction":"rises"},{"sphere_name":"higher","constant_capital":92,"variable_capital":8,"wage_factor":125,"composition":"higher","old_price":12000,"new_price":11657,"price_change":-343,"direction":"falls"}]',
    '2026-05-01 00:00:00.000000'
);

-- +goose Down
DELETE FROM wage_effect_analyses WHERE id IN (
    '5eed000000000000001101'
);
