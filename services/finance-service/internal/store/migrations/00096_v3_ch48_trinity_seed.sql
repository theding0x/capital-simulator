-- +goose Up
-- Ch. 48 — Trinity Formula: three revenue streams (capital/land/labour), one
-- formula binding them, three fetish forms. Average profit 20 + ground-rent 3 =
-- surplus-value 23; wages 80 = variable capital. Apparent revenue 103 = s + v.
INSERT INTO `revenue_streams` (`id`,`source`,`apparent_revenue_bp`,`actual_source_bp`,`is_fetishised`,`created_at`) VALUES
    ('5eed000000000000004801', 1, 20, 20, 1, '1894-01-01 00:00:00.000000'),
    ('5eed000000000000004802', 2, 3, 3, 1, '1894-01-01 00:00:00.000000'),
    ('5eed000000000000004803', 3, 80, 0, 1, '1894-01-01 00:00:00.000000');
INSERT INTO `trinity_formulas` (`id`,`capital_stream_id`,`land_stream_id`,`labour_stream_id`,`total_surplus_value_bp`,`total_apparent_revenue_bp`,`created_at`) VALUES
    ('5eed000000000000004800', '5eed000000000000004801', '5eed000000000000004802', '5eed000000000000004803', 23, 103, '1894-01-01 00:00:00.000000');
INSERT INTO `revenue_fetish_forms` (`id`,`source`,`surface_formula`,`real_relation`,`mystification_kind`,`created_at`) VALUES
    ('5eed000000000000004804', 1, 'Capital — Interest', 'Capital extracts surplus-labour from wage-workers; interest is one share of that surplus-value', 'value appears to beget value of itself (100 becomes 110)', '1894-01-01 00:00:00.000000'),
    ('5eed000000000000004805', 2, 'Land — Ground-Rent', 'Landed property captures a portion of surplus-value by barring free access to the soil', 'a use-value without value (land) is paired with an exchange-value (rent)', '1894-01-01 00:00:00.000000'),
    ('5eed000000000000004806', 3, 'Labour — Wages', 'Wages are the price of labour-power and recover variable capital, not a payment for labour itself', 'labour, the very substance of value, appears as a commodity with its own price', '1894-01-01 00:00:00.000000');

-- +goose Down
DELETE FROM `revenue_fetish_forms` WHERE `id` IN ('5eed000000000000004804','5eed000000000000004805','5eed000000000000004806');
DELETE FROM `trinity_formulas` WHERE `id` = '5eed000000000000004800';
DELETE FROM `revenue_streams` WHERE `id` IN ('5eed000000000000004801','5eed000000000000004802','5eed000000000000004803');
