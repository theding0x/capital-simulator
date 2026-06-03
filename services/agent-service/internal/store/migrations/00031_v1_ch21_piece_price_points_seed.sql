-- +goose Up
-- Seed the Ch. 21 dynamic piece-price series (#219) for the Spinning Mill
-- worker (agent 5eed…1900ff), whose piece-wage seed is 6 farthings at 24
-- pieces a day = 144 farthings (3s.) implied daily wage.
--
-- The series demonstrates Marx's law of piece-wages: as productivity rises
-- 1.0 -> 1.5 -> 2.0 the implied daily wage stays 144 farthings while the price
-- per piece falls 6f -> 4f -> 3f and the normal output rises 24 -> 36 -> 48.
-- At 2x productivity the piece price is exactly halved. The dashboard comes up
-- showing this falling-price-with-rising-productivity curve.
INSERT INTO piece_price_points
    (id, agent_id, sequence_num, productivity_factor, price_pence, normal_output, occurred_at)
VALUES
    ('5eed00000000000000210a01', '5eed000000000000001900ff', 0, 1.0, 6, 24, '2026-05-14 10:00:00.000000'),
    ('5eed00000000000000210a02', '5eed000000000000001900ff', 1, 1.5, 4, 36, '2026-05-14 11:00:00.000000'),
    ('5eed00000000000000210a03', '5eed000000000000001900ff', 2, 2.0, 3, 48, '2026-05-14 12:00:00.000000');

-- +goose Down
DELETE FROM piece_price_points
WHERE id IN (
    '5eed00000000000000210a01',
    '5eed00000000000000210a02',
    '5eed00000000000000210a03'
);
