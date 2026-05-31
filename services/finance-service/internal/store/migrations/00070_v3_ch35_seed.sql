-- +goose Up
INSERT INTO gold_reserves (id, name, amount, description, created_at) VALUES
('5eed000000000000003500', 'BoE Gold Reserve Nov 1857', 800000000, 'Bank of England gold reserve in November 1857 during the crisis that led to suspension of the Bank Act 1844.', '1857-11-01 00:00:00.000000'),
('5eed000000000000003501', 'BoE Gold Reserve 1844 Act Baseline', 1400000000, 'Bank of England gold reserve at the time of the Bank Charter Act 1844, forming the baseline for the fiduciary issue calculation.', '1844-07-19 00:00:00.000000');

INSERT INTO rates_of_exchange (id, name, deviation_bp, description, created_at) VALUES
('5eed000000000000003510', 'London-Hamburg 1847 Unfavourable', -150, 'London-Hamburg exchange rate in October 1847: adverse balance of payments driven by corn imports and railway speculation produced an unfavourable rate of -150 basis points below par.', '1847-10-15 00:00:00.000000'),
('5eed000000000000003511', 'London-India 1850s Adverse Drain', -200, 'London-India exchange rate in the 1850s: persistent adverse balance from unilateral tribute payments and commodity imports produced a sustained gold drain of -200 basis points below par.', '1857-06-01 00:00:00.000000');

-- +goose Down
DELETE FROM rates_of_exchange WHERE id IN ('5eed000000000000003510', '5eed000000000000003511');
DELETE FROM gold_reserves WHERE id IN ('5eed000000000000003500', '5eed000000000000003501');
