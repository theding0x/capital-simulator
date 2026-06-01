-- +goose Up
-- Ch. 50 — three competition illusions, one cost-price fetish (c=70/v=30 → cost 100,
-- profit 20, actual value 120), one society-wide annual revenue account
-- (c=6000, v=1500, s=1500 → new value 3000, total 9000), two value-appearance contrasts.
INSERT INTO `competition_illusions` (`id`,`name`,`surface_appearance`,`real_relation`,`produced_by`,`created_at`) VALUES
    ('5eed000000000000005000', 'Profit arises in exchange', 'profit appears as a markup added above cost-price in the sale', 'profit is surplus-value created in production by unpaid labour', 'price_competition', '1894-01-01 00:00:00.000000'),
    ('5eed000000000000005001', 'Interest arises from capital as a thing', 'money appears to breed money of itself, independent of labour', 'interest is a share of surplus-value extracted from wage-workers', 'interest_form', '1894-01-01 00:00:00.000000'),
    ('5eed000000000000005002', 'Wages pay for labour', 'wages appear to buy labour at its full value', 'wages buy labour-power, whose use creates more value than it costs', 'wage_form', '1894-01-01 00:00:00.000000');
INSERT INTO `cost_price_fetishes` (`id`,`constant_capital_bp`,`variable_capital_bp`,`cost_price_bp`,`profit_bp`,`actual_value_bp`,`created_at`) VALUES
    ('5eed000000000000005003', 70, 30, 100, 20, 120, '1894-01-01 00:00:00.000000');
INSERT INTO `annual_revenue_accounts` (`id`,`new_value_created_lm`,`wages_share_lm`,`profit_and_rent_share_lm`,`constant_capital_lm`,`total_social_product_lm`,`created_at`) VALUES
    ('5eed000000000000005004', 3000, 1500, 1500, 6000, 9000, '1894-01-01 00:00:00.000000');
INSERT INTO `value_appearance_contrasts` (`id`,`illusion_id`,`surface`,`reality`,`created_at`) VALUES
    ('5eed000000000000005005', '5eed000000000000005000', 'the capitalist sells above cost and pockets the difference', 'the difference is unpaid surplus-labour, not a gain made in circulation', '1894-01-01 00:00:00.000000'),
    ('5eed000000000000005006', '5eed000000000000005002', 'a fair day''s wage for a fair day''s work', 'the working day exceeds the labour-time the wage represents', '1894-01-01 00:00:00.000000');

-- +goose Down
DELETE FROM `value_appearance_contrasts` WHERE `id` IN ('5eed000000000000005005','5eed000000000000005006');
DELETE FROM `annual_revenue_accounts` WHERE `id` = '5eed000000000000005004';
DELETE FROM `cost_price_fetishes` WHERE `id` = '5eed000000000000005003';
DELETE FROM `competition_illusions` WHERE `id` IN ('5eed000000000000005000','5eed000000000000005001','5eed000000000000005002');
