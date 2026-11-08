-- +goose Up
-- Corrects the Ch. 7 §2 spinning seed from 00007 (issue #450).
-- Scale: 1s = 120 labour-minutes.
--
-- 00007 ran a 12-hour (720 min) day on the 6-hour means (10 lbs cotton, spindle
-- 2s), giving a product of 18s. Marx: a 12-hour day absorbs 20 lbs cotton (20s)
-- and 4s of spindle, adds 6s, so 20 lbs yarn = 30s; the capitalist advances
-- 24s + 3s (labour-power) = 27s and recovers 30s, surplus 3s.
--   means    = 20 x 120 + 480 = 2,880 min (24s)
--   advanced = 2,880 + 360    = 3,240 min (27s)
UPDATE labour_processes
    SET means_json = '{"raw_materials":[{"commodity_id":"cotton","quantity":20,"snlt_per_unit":120}],"instruments":[{"commodity_id":"spindle","wear_per_run":480}]}',
        product_quantity = 20
    WHERE id = '5eed00000000000000070001';
UPDATE labour_capitalists SET money_capital = 3240
    WHERE id = '5eed00000000000000060010';

-- +goose Down
UPDATE labour_processes
    SET means_json = '{"raw_materials":[{"commodity_id":"cotton","quantity":10,"snlt_per_unit":120}],"instruments":[{"commodity_id":"spindle","wear_per_run":240}]}',
        product_quantity = 10
    WHERE id = '5eed00000000000000070001';
UPDATE labour_capitalists SET money_capital = 1800
    WHERE id = '5eed00000000000000060010';
