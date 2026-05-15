-- +goose Up
-- § core example: "capital of £1,000 begets yearly a surplus-value of £200"
-- c=800, v=200, s'=100%, repayment period = 5 years
INSERT INTO reproduction_cycles (id, constant_capital, variable_capital, surplus_rate, periods, created_at)
VALUES ('5eed000000000000002301', 800, 200, 1.0, 5, NOW(6));

-- +goose Down
DELETE FROM reproduction_cycles WHERE id IN ('5eed000000000000002301');
