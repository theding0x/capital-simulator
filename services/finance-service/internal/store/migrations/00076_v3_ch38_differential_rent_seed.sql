-- +goose Up
INSERT INTO monopolised_natural_forces (id, kind, is_monopolisable, parcel_id, created_at) VALUES
    ('5eed000000000000003800', 'waterfall', 1, '5eed000000000000003700', '1867-01-01 00:00:00.000000');
INSERT INTO differential_rents_ch38 (id, parcel_id, surplus_profit_bp, annual_rent_labour_minutes, created_at) VALUES
    ('5eed000000000000003801', '5eed000000000000003700', 2500, 4800, '1867-01-01 00:00:00.000000');
INSERT INTO capitalised_rent_prices (id, parcel_id, annual_rent_labour_minutes, interest_rate_bp, capitalised_price_labour_minutes, created_at) VALUES
    ('5eed000000000000003802', '5eed000000000000003700', 4800, 500, 96000, '1867-01-01 00:00:00.000000');
INSERT INTO pop_surplus_profits (id, capital_id, general_price_of_production_bp, individual_price_of_production_bp, surplus_profit_bp, created_at) VALUES
    ('5eed000000000000003803', 'waterfall-factory-1867', 11500, 9000, 2500, '1867-01-01 00:00:00.000000');

-- +goose Down
DELETE FROM pop_surplus_profits WHERE id = '5eed000000000000003803';
DELETE FROM capitalised_rent_prices WHERE id = '5eed000000000000003802';
DELETE FROM differential_rents_ch38 WHERE id = '5eed000000000000003801';
DELETE FROM monopolised_natural_forces WHERE id = '5eed000000000000003800';
