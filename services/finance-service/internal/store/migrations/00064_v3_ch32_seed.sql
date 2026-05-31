-- +goose Up
-- Vol. III Ch. 32 seeds — capital released onto the loan market without real
-- expansion behind it. Seed IDs start with 5eed000000000000003200..3201 (chapter 32).
--
-- Seed 3200: Capital released by a fall in raw-material prices (cause=1), NOT
--            reinvested. The same reproduction now ties up less money; the
--            difference is lent out. It merely OVERSTATES real accumulation —
--            overstatement_pct = 2500 bp (25%).
--
-- Seed 3201: The growing rentier class (cause=3) — savings withdrawn from
--            productive investment and placed at interest, swelling loan capital
--            with no new productive capital. NOT reinvested; overstatement_pct =
--            10000 bp (the whole amount is apparent, not real, accumulation).

INSERT INTO `capital_releases`
    (`id`, `name`, `released_amount`, `cause`, `reinvested`, `overstatement_pct`, `overstatement_description`, `description`, `created_at`)
VALUES
    ('5eed000000000000003200',
     'Capital released by raw-material price fall',
     500000000,
     1,
     FALSE,
     2500,
     'freed money lent out but not reinvested; loan capital overstates real accumulation',
     'raw-material prices fall; the same reproduction ties up less money; the difference is lent out',
     '2025-01-01 00:00:00.000000'),
    ('5eed000000000000003201',
     'Rentier-class savings placed at interest',
     200000000,
     3,
     FALSE,
     10000,
     'the whole sum is apparent, not real, accumulation — withdrawn from productive investment',
     'the growing rentier class lives off interest, swelling loan capital independently of real expansion',
     '2025-01-01 00:01:00.000000');

-- +goose Down
DELETE FROM `capital_releases`
WHERE `id` IN ('5eed000000000000003200', '5eed000000000000003201');
