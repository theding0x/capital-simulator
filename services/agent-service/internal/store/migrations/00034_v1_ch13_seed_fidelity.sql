-- +goose Up
-- Corrects the Ch. 13 §3 simple co-operation seed from 00011 (issue #450).
-- 00011 staged the 12-man fixture as Arkwright's spinning floor, which is a
-- Ch. 15 factory (machinery). Marx's simple co-operation example is the
-- "12 masons, in their collective working-day of 144 hours" building together,
-- so the capitalist, the ten workers and the co-operation are renamed.
-- Overlookers A and B stay: Marx's overlookers appear in section 4.
-- Sizes, working days and supervisory flags are unchanged.
UPDATE agents SET name = 'Master Builder' WHERE id = '5eed00000000000000130100';
UPDATE agents SET name = 'Mason 01' WHERE id = '5eed00000000000000130001';
UPDATE agents SET name = 'Mason 02' WHERE id = '5eed00000000000000130002';
UPDATE agents SET name = 'Mason 03' WHERE id = '5eed00000000000000130003';
UPDATE agents SET name = 'Mason 04' WHERE id = '5eed00000000000000130004';
UPDATE agents SET name = 'Mason 05' WHERE id = '5eed00000000000000130005';
UPDATE agents SET name = 'Mason 06' WHERE id = '5eed00000000000000130006';
UPDATE agents SET name = 'Mason 07' WHERE id = '5eed00000000000000130007';
UPDATE agents SET name = 'Mason 08' WHERE id = '5eed00000000000000130008';
UPDATE agents SET name = 'Mason 09' WHERE id = '5eed00000000000000130009';
UPDATE agents SET name = 'Mason 10' WHERE id = '5eed00000000000000130010';
UPDATE cooperations SET name = 'Twelve Masons' WHERE id = '5eed00000000000000130c01';

-- +goose Down
UPDATE agents SET name = 'Richard Arkwright' WHERE id = '5eed00000000000000130100';
UPDATE agents SET name = 'Spinner 01' WHERE id = '5eed00000000000000130001';
UPDATE agents SET name = 'Spinner 02' WHERE id = '5eed00000000000000130002';
UPDATE agents SET name = 'Spinner 03' WHERE id = '5eed00000000000000130003';
UPDATE agents SET name = 'Spinner 04' WHERE id = '5eed00000000000000130004';
UPDATE agents SET name = 'Spinner 05' WHERE id = '5eed00000000000000130005';
UPDATE agents SET name = 'Spinner 06' WHERE id = '5eed00000000000000130006';
UPDATE agents SET name = 'Spinner 07' WHERE id = '5eed00000000000000130007';
UPDATE agents SET name = 'Spinner 08' WHERE id = '5eed00000000000000130008';
UPDATE agents SET name = 'Spinner 09' WHERE id = '5eed00000000000000130009';
UPDATE agents SET name = 'Spinner 10' WHERE id = '5eed00000000000000130010';
UPDATE cooperations SET name = 'Arkwright Spinning Floor' WHERE id = '5eed00000000000000130c01';
