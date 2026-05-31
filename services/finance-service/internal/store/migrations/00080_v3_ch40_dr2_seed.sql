-- +goose Up
-- Marx Ch. 40 — doubled successive investment on Soil D (Norfolk Grade-D parcel
-- 5eed…3700). Each tranche mirrors the DR I Grade-D row: £2.50 (250), 4 qtr (400),
-- surplus-profit £9 (900). Two tranches → TotalRent 1800 = £18. Lease rent locked
-- at 900 (single-investment level) so the capitalist pockets the DR II uplift.
INSERT INTO `successive_investments`
    (`id`, `parcel_id`, `tranche_number`, `capital_bp`, `output_quarters`, `surplus_profit_bp`, `created_at`)
VALUES
    ('5eed000000000000004000', '5eed000000000000003700', 1, 250, 400, 900, '1867-01-01 00:00:00.000000'),
    ('5eed000000000000004001', '5eed000000000000003700', 2, 250, 400, 900, '1867-01-01 00:00:00.000000');

INSERT INTO `lease_terms`
    (`id`, `parcel_id`, `start_year`, `duration_years`, `fixed_rent_bp`, `created_at`)
VALUES
    ('5eed000000000000004002', '5eed000000000000003700', 1867, 21, 900, '1867-01-01 00:00:00.000000');

-- +goose Down
DELETE FROM `successive_investments` WHERE `id` IN ('5eed000000000000004000', '5eed000000000000004001');
DELETE FROM `lease_terms`            WHERE `id` = '5eed000000000000004002';
