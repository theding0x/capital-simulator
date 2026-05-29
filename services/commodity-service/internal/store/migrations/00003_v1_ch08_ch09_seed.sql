-- +goose Up
-- Seed data for Capital Vol. I, Ch. 8 (Constant and Variable Capital) and
-- Ch. 9 (The Rate of Surplus-Value).
--
-- Each row is one production run: c (constant capital, value transferred from
-- means of production), v (variable capital, value of labour-power), s
-- (surplus-value, the unpaid portion of the working day). All magnitudes are
-- in labour-minutes; rates are derived as s/v at read time.

INSERT IGNORE INTO production_accounts
    (id, constant, variable, surplus, created_at)
VALUES
    -- Ch.7 §2 spinning mill (carried forward from the labour_processes seed):
    -- c = 1440 (10 lbs cotton @ 120 + spindle 240); v = 360; s = 360.
    -- s/v = 100%, the canonical Marx example.
    ('5eed0000000000000008c701', 1440, 360, 360, '2026-05-01 10:00:00.000000'),

    -- Ch.9 §1 worked example: "C = £410 const. + £90 var.; surplus-value = £90".
    -- 1 shilling = 120 min, 1 penny = 10 min. £410 = 8200s = 984000d → 9_840_000 min.
    -- £90 = 1800s = 21600d → 216_000 min. Same conversion for surplus.
    -- s/v = 100%; the chapter's headline arithmetic.
    ('5eed0000000000000009a101', 9840000, 216000, 216000, '2026-05-01 12:00:00.000000'),

    -- Ch.9 §1 spinning mill (1871): £52 var. + £80 surpl. (no constant given;
    -- modelled here as a typical 4× organic composition c = 4v).
    -- £52 = 1040s = 12480d → 124_800 min. £80 = 1600s = 19200d → 192_000 min.
    -- Constant chosen as 499_200 (= 4 × 124_800) for a plausible c:v of 4:1.
    -- s/v ≈ 153.85%, the rate Marx cites as "153 11/13 per cent".
    ('5eed0000000000000009a102', 499200, 124800, 192000, '2026-05-01 13:00:00.000000'),

    -- Ch.9 §1 Jacob's wheat (1815): wages £3 10s. (= 70s = 8400d → 84_000 min);
    -- surplus £3 11s. (= 71s = 8520d → 85_200 min). Constant left at 0 since
    -- Marx's table lumps seed/rent into the same column; rate s/v slightly above 1.
    ('5eed0000000000000009a103', 0, 84000, 85200, '2026-05-01 14:00:00.000000');

-- +goose Down
DELETE FROM production_accounts WHERE id IN (
    '5eed0000000000000008c701',
    '5eed0000000000000009a101',
    '5eed0000000000000009a102',
    '5eed0000000000000009a103'
);
