-- +goose Up
-- The Ch. 7 spinning-mill production account (00003) carries forward the
-- 12-hour run, so c must be the 12-hour means: 20 lbs cotton + 4s spindle =
-- 24s = 2,880 min (1s = 120 min), not the 6-hour 1,440 (issue #450).
UPDATE production_accounts SET constant = 2880
    WHERE id = '5eed0000000000000008c701';

-- +goose Down
UPDATE production_accounts SET constant = 1440
    WHERE id = '5eed0000000000000008c701';
