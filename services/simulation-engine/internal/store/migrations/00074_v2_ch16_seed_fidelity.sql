-- +goose Up
-- Vol. II Ch. 16 seed fidelity (issue #450, v2-ch16 item).
--
-- Each seeded magnitude below is tagged M (Marx's own figure), D (derived from his
-- figures by the arithmetic shown) or I (illustrative: assumed, not Marx's). 00051 is
-- append-only and is not edited. This migration restates two of its four fixtures in
-- place, by seed id: no rows are added or removed, so Down needs no DELETE.
--
-- The year. Marx's Ch. 16 year is "a year of 50 weeks" (sec. I, and Ch. 15 likewise), so
-- 50 x 7 x 1,440 = 504,000 minutes and n = 504,000 / turnover_period_minutes, stored as
-- round-half-up(10,000 x n) basis points (D). 00051 applied it to Capitals A and B only:
-- its spinner used a 52-week n and its locomotive works an n (43,461) that fits no year.
--
-- Left as seeded. Capital A (ids ...1601, 1607, 1611) and Capital B (...1602, 1608,
-- 1612) are Marx's own (M), Ch. 16 sec. I. A: 500 pounds, a 5-week period, n = 10,
-- s' = 100%, annual rate 1,000%. B: 5,000 pounds, a 50-week period, n = 1, s' = 100%,
-- annual rate 100%. Their contrast (...1615) is 10 to 1 (D). 00051's header quotes
-- for A and B are paraphrases and put them in "sec. 2", which is sec. I of Ch. 16, but
-- the figures are right.
--
-- Spinning Mill (ids ...1603, 1609, 1613). The fixture is Vol. II Ch. 1 (reused in
-- Ch. 2 and 3), not Ch. 16.
--   v   12,000p = 50 pounds, the weekly wage of the mill's 50 labourers      (M, Ch. 1)
--   s   18,720p = 78 pounds of surplus-value on a week's yarn                (M, Ch. 1)
--   s/v 15,600 bp: 10,000 x 18,720 / 12,000 exactly, "156 per cent"          (M, Ch. 1)
--   period 10,080 min = one week, unchanged from 00051                       (I: Ch. 1 gives
--          the mill's wage and surplus per week, but taking the week as its turnover
--          period is assumed)
--   n   500,000 bp = 504,000 / 10,080 = 50 turnovers                         (D)
--   annual 780,000 bp = 500,000 x 15,600 / 10,000 = 7,800%                   (D)
-- 00051 had 100 and 156 pounds (twice Marx's, with the ratio right), n = 52 and an annual
-- rate of 8,112%. Marx gives no annual rate for the mill. Over the year the mill pays
-- 50 x 50 = 2,500 pounds in wages and yields 50 x 78 = 3,900 pounds of surplus (D).
-- Ch. 1's own passage also gives 1,500 of the week's 3,000 hours as surplus labour (100%)
-- beside its "156 per cent". The seed follows the stated 156%, which the locomotive
-- works inherits.
--
-- Locomotive Works (ids ...1605, 1610, 1614). Marx gives it no rate of surplus-value,
-- advance or turnover number, in Ch. 16 or elsewhere in Vol. II. Every value below
-- except the period is derived from statements in Ch. 12, not stated by Marx (the period
-- is his twelve weeks). Ch. 12 gives a three-month working period, the weekly expenditure
-- "multiplied by twelve", payment on delivery of the finished locomotive (so no
-- circulation time), and the same division into necessary and surplus labour as the
-- spinner. So the seed is the spinner's week laid out for twelve weeks:
--   period 120,960 min = 12 x 10,080, unchanged from 00051                   (M: Ch. 12, twelve weeks)
--   v   144,000p = 12 x 12,000p = 600 pounds                                 (D)
--   s   224,640p = 12 x 18,720p = 936 pounds                                 (D)
--   s/v 15,600 bp: the spinner's rate, carried over by Ch. 12's "same" division (D)
--   n   41,667 bp = 504,000 / 120,960 = 4.1667, 41,666.67 rounded half up    (D)
--   annual 65,001 bp = 41,667 x 15,600 / 10,000 = 65,000.52, rounded half up (D)
-- The annual rate is 650% (156 x 50 / 12 = 650) to within the 1 bp that rounding n costs,
-- and the spinner's is twelve times it. 00051's 800 pounds, 80% and n = 43,461 had no
-- source (I), and the n fits no year: 50 weeks over twelve gives 41,667, 52 weeks 43,333
-- and 365 days 43,452.
--
-- Left illustrative (I): the period label '1871' (borrowed from the spinning-mill fixture
-- for every row), the three industrial_capital_id strings and created_at.
--
-- Known limitation: RecordPerCircuitRate upserts with id = VALUES(id), so on a database
-- where someone re-recorded a seeded advance's rate the seeded per-circuit id no longer
-- exists and that statement matches no row. Fresh volumes are unaffected. The statements
-- address rows by seed id on purpose, so a user's own rows are never touched.

UPDATE variable_capital_advances SET pence = 12000, turnovers_per_year_basis_points = 500000 WHERE id = '5eed0000000000001603';
UPDATE per_circuit_surplus_rates SET surplus_pence = 18720, variable_pence = 12000, basis_points = 15600 WHERE id = '5eed0000000000001609';
UPDATE annual_surplus_rates SET n_basis_points = 500000, per_circuit_basis_points = 15600, annual_basis_points = 780000 WHERE id = '5eed0000000000001613';

UPDATE variable_capital_advances SET pence = 144000, turnovers_per_year_basis_points = 41667 WHERE id = '5eed0000000000001605';
UPDATE per_circuit_surplus_rates SET surplus_pence = 224640, variable_pence = 144000, basis_points = 15600 WHERE id = '5eed0000000000001610';
UPDATE annual_surplus_rates SET n_basis_points = 41667, per_circuit_basis_points = 15600, annual_basis_points = 65001 WHERE id = '5eed0000000000001614';

-- +goose Down
-- Restores 00051 exactly: the 100-pound spinner at n = 52 (8,112%) and the 800-pound,
-- 80%, n = 43,461 locomotive works.
UPDATE annual_surplus_rates SET n_basis_points = 43461, per_circuit_basis_points = 8000, annual_basis_points = 34769 WHERE id = '5eed0000000000001614';
UPDATE per_circuit_surplus_rates SET surplus_pence = 153600, variable_pence = 192000, basis_points = 8000 WHERE id = '5eed0000000000001610';
UPDATE variable_capital_advances SET pence = 192000, turnovers_per_year_basis_points = 43461 WHERE id = '5eed0000000000001605';

UPDATE annual_surplus_rates SET n_basis_points = 520000, per_circuit_basis_points = 15600, annual_basis_points = 811200 WHERE id = '5eed0000000000001613';
UPDATE per_circuit_surplus_rates SET surplus_pence = 37440, variable_pence = 24000, basis_points = 15600 WHERE id = '5eed0000000000001609';
UPDATE variable_capital_advances SET pence = 24000, turnovers_per_year_basis_points = 520000 WHERE id = '5eed0000000000001603';
