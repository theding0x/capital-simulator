-- +goose Up
-- Vol. II Ch. 13 — corrective seed for the Ch. 13 working periods (5eed…1301-1306).
-- Marx: "The time of production is here longer than the labour-time." The
-- working period is the connected labour-time; the production time contains it.
-- 00047 had set working_day_count to the whole production span. Here the working
-- period becomes the labour-days (10-hour working-days, working_day_length_minutes
-- = 600) and each gap's labour_time_nanos is tied to it: days × 600 min × 6e10 ns.
-- Production spans are raised to the Ch. 12 natural floors (corn 365, cattle 1095).
-- 1 calendar day = 86400000000000 ns.

UPDATE working_periods SET working_day_count = 10 WHERE id = '5eed000000000000001301';
UPDATE working_periods SET working_day_count = 10 WHERE id = '5eed000000000000001302';
UPDATE working_periods SET working_day_count = 5  WHERE id = '5eed000000000000001303';
UPDATE working_periods SET working_day_count = 30 WHERE id = '5eed000000000000001304';
UPDATE working_periods SET working_day_count = 3  WHERE id = '5eed000000000000001305';
UPDATE working_periods SET working_day_count = 5  WHERE id = '5eed000000000000001306';

-- Corn:      10 labour-days /   365 production-days (Ch. 12 floor 365).
-- Wine:      10 labour-days /   365 production-days.
-- Cattle:     5 labour-days /  1095 production-days (Ch. 12 floor 1095).
-- Timber:    30 labour-days / 18250 production-days.
-- Porcelain:  3 labour-days /    30 production-days.
-- Fallow:     5 labour-days /   180 production-days.
UPDATE production_labour_gaps SET production_time_nanos = 31536000000000000,   labour_time_nanos = 360000000000000,   gap_nanos = 31176000000000000   WHERE id = '5eed000000000000001310';
UPDATE production_labour_gaps SET production_time_nanos = 31536000000000000,   labour_time_nanos = 360000000000000,   gap_nanos = 31176000000000000   WHERE id = '5eed000000000000001311';
UPDATE production_labour_gaps SET production_time_nanos = 94608000000000000,   labour_time_nanos = 180000000000000,   gap_nanos = 94428000000000000   WHERE id = '5eed000000000000001312';
UPDATE production_labour_gaps SET production_time_nanos = 1576800000000000000, labour_time_nanos = 1080000000000000,  gap_nanos = 1575720000000000000 WHERE id = '5eed000000000000001313';
UPDATE production_labour_gaps SET production_time_nanos = 2592000000000000,    labour_time_nanos = 108000000000000,    gap_nanos = 2484000000000000    WHERE id = '5eed000000000000001314';
UPDATE production_labour_gaps SET production_time_nanos = 15552000000000000,   labour_time_nanos = 180000000000000,   gap_nanos = 15372000000000000   WHERE id = '5eed000000000000001315';

-- +goose Down
UPDATE working_periods SET working_day_count = 280   WHERE id = '5eed000000000000001301';
UPDATE working_periods SET working_day_count = 365   WHERE id = '5eed000000000000001302';
UPDATE working_periods SET working_day_count = 730   WHERE id = '5eed000000000000001303';
UPDATE working_periods SET working_day_count = 18250 WHERE id = '5eed000000000000001304';
UPDATE working_periods SET working_day_count = 30    WHERE id = '5eed000000000000001305';
UPDATE working_periods SET working_day_count = 180   WHERE id = '5eed000000000000001306';

UPDATE production_labour_gaps SET production_time_nanos = 24192000000000000,   labour_time_nanos = 864000000000000,  gap_nanos = 23328000000000000   WHERE id = '5eed000000000000001310';
UPDATE production_labour_gaps SET production_time_nanos = 31536000000000000,   labour_time_nanos = 864000000000000,  gap_nanos = 30672000000000000   WHERE id = '5eed000000000000001311';
UPDATE production_labour_gaps SET production_time_nanos = 63072000000000000,   labour_time_nanos = 432000000000000,  gap_nanos = 62640000000000000   WHERE id = '5eed000000000000001312';
UPDATE production_labour_gaps SET production_time_nanos = 1576800000000000000, labour_time_nanos = 2592000000000000, gap_nanos = 1574208000000000000 WHERE id = '5eed000000000000001313';
UPDATE production_labour_gaps SET production_time_nanos = 2592000000000000,    labour_time_nanos = 259200000000000,  gap_nanos = 2332800000000000    WHERE id = '5eed000000000000001314';
UPDATE production_labour_gaps SET production_time_nanos = 15552000000000000,   labour_time_nanos = 432000000000000,  gap_nanos = 15120000000000000   WHERE id = '5eed000000000000001315';
