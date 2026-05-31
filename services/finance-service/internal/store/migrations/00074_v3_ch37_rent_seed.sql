-- +goose Up
INSERT INTO landowners (id, name, created_at) VALUES
    ('5eed000000000000003701', 'Duke of Sutherland', '1867-01-01 00:00:00.000000');

INSERT INTO land_parcels (id, owner_id, fertility_grade, area_hectares, location, created_at) VALUES
    ('5eed000000000000003700', '5eed000000000000003701', 4, 100, 'Norfolk capitalist tenant-farm (Grade D)', '1867-01-01 00:00:00.000000');

INSERT INTO agricultural_capitalists (id, capital_advanced, lease_parcel_id, created_at) VALUES
    ('5eed000000000000003702', 4800, '5eed000000000000003700', '1867-01-01 00:00:00.000000');

INSERT INTO ground_rents (id, parcel_id, capitalist_id, land_owner_id, `form`, amount_labour_minutes, period_years, created_at) VALUES
    ('5eed000000000000003703', '5eed000000000000003700', '5eed000000000000003702', '5eed000000000000003701', 1, 2400, 1, '1867-01-01 00:00:00.000000');

-- +goose Down
DELETE FROM ground_rents WHERE id = '5eed000000000000003703';
DELETE FROM agricultural_capitalists WHERE id = '5eed000000000000003702';
DELETE FROM land_parcels WHERE id = '5eed000000000000003700';
DELETE FROM landowners WHERE id = '5eed000000000000003701';
