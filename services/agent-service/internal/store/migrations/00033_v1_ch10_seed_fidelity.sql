-- +goose Up
-- Corrects the Ch. 10 §6 Factory Act (1847, Ten Hours' Bill) working-day seed
-- from 00009 (issue #450). 00009 split the 10-hour day 5h + 5h, but the chapter
-- holds necessary labour at 6 hours (3s of labour-power, 1s = 120 min), so the
-- statutory 600-minute day is 360 necessary + 240 surplus (s' = 2/3).
UPDATE working_days
    SET necessary_labour_minutes = 360, surplus_labour_minutes = 240
    WHERE id = '5eed00000000000000100004';

-- +goose Down
UPDATE working_days
    SET necessary_labour_minutes = 300, surplus_labour_minutes = 300
    WHERE id = '5eed00000000000000100004';
