-- +goose Up
-- Corrects two Ch. 9 seeds from 00003 so they trace to Marx's text
-- (issue #450). Scale: 1s = 120 labour-minutes, so 1 pound = 2,400.
--
--   * Ch.9 §1 worked example, "C = £410 const. + £90 var.":
--     £410 = 410 x 2,400 = 984,000 min. 00003 seeded 9,840,000 (10x too
--     large), giving c:v of 45.6 instead of Marx's 4.56.
--   * Ch.9 §1 spinning mill (1871): Marx states "the constant portion of the
--     value of the week's product is £378", so c = 378 x 2,400 = 907,200 min.
--     00003 invented c = 4v = 499,200. v (£52) and s (£80) were already right.
UPDATE production_accounts SET constant = 984000
    WHERE id = '5eed0000000000000009a101';
UPDATE production_accounts SET constant = 907200
    WHERE id = '5eed0000000000000009a102';

-- +goose Down
UPDATE production_accounts SET constant = 9840000
    WHERE id = '5eed0000000000000009a101';
UPDATE production_accounts SET constant = 499200
    WHERE id = '5eed0000000000000009a102';
