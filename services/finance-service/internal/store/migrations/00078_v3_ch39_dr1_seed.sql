-- +goose Up
INSERT INTO dr1_tables (id, description, total_rental_bp, regulating_grade, entries_json, created_at) VALUES (
    '5eed000000000000003900',
    'Marx Table I (Ch. 39): four soils A-D, equal capital £2.50, regulating price £3/qtr',
    1800,
    1,
    '[{"id":"5eed000000000000003901","table_id":"5eed000000000000003900","grade":1,"capital_advanced_bp":250,"output_quarters":100,"price_of_production_bp":300,"surplus_product_qtrs":0,"money_rent_bp":0,"created_at":"1867-01-01T00:00:00Z"},{"id":"5eed000000000000003902","table_id":"5eed000000000000003900","grade":2,"capital_advanced_bp":250,"output_quarters":200,"price_of_production_bp":300,"surplus_product_qtrs":100,"money_rent_bp":300,"created_at":"1867-01-01T00:00:00Z"},{"id":"5eed000000000000003903","table_id":"5eed000000000000003900","grade":3,"capital_advanced_bp":250,"output_quarters":300,"price_of_production_bp":300,"surplus_product_qtrs":200,"money_rent_bp":600,"created_at":"1867-01-01T00:00:00Z"},{"id":"5eed000000000000003904","table_id":"5eed000000000000003900","grade":4,"capital_advanced_bp":250,"output_quarters":400,"price_of_production_bp":300,"surplus_product_qtrs":300,"money_rent_bp":900,"created_at":"1867-01-01T00:00:00Z"}]',
    '1867-01-01 00:00:00.000000'
);

-- +goose Down
DELETE FROM dr1_tables WHERE id = '5eed000000000000003900';
