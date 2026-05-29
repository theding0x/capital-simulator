-- +goose Up
-- § spinner example §1: £10,000 capital (£8,000c + £2,000v), s'=100%, full accumulation,
-- composition 4/5 constant + 1/5 variable, 3 periods — the "Abraham begat Isaac" sequence.
INSERT INTO accumulation_scenarios (id, name, constant_capital, variable_capital, surplus_rate, accum_rate, composition_ratio, periods, created_at)
VALUES ('5eed000000000000002401', 'Spinner — §1 full accumulation (Abraham)', 8000, 2000, 1.0, 1.0, 0.8, 3, NOW(6));

-- § partial accumulation §3: capitalist consumes half; revenue = £1,000.
INSERT INTO accumulation_scenarios (id, name, constant_capital, variable_capital, surplus_rate, accum_rate, composition_ratio, periods, created_at)
VALUES ('5eed000000000000002402', 'Spinner — §3 partial accumulation (rate=0.5)', 8000, 2000, 1.0, 0.5, 0.8, 3, NOW(6));

-- +goose Down
DELETE FROM accumulation_scenarios WHERE id IN (
    '5eed000000000000002401',
    '5eed000000000000002402'
);
