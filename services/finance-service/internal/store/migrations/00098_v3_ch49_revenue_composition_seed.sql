-- +goose Up
-- Ch. 49 — annual product c=6000, v=1500, s=1500 → total 9000, new value 3000;
-- revenue limit max=3000 (v+s), constant replacement 6000; partition s=1500 →
-- profit of enterprise 900 + interest 300 + ground-rent 300.
INSERT INTO `annual_value_compositions` (`id`,`constant_capital_labour_minutes`,`variable_capital_labour_minutes`,`surplus_value_labour_minutes`,`total_value_labour_minutes`,`new_value_labour_minutes`,`created_at`) VALUES
    ('5eed000000000000004900', 6000, 1500, 1500, 9000, 3000, '1894-01-01 00:00:00.000000');
INSERT INTO `revenue_limits` (`id`,`composition_id`,`max_distributable_revenue_lm`,`constant_capital_replacement`,`created_at`) VALUES
    ('5eed000000000000004901', '5eed000000000000004900', 3000, 6000, '1894-01-01 00:00:00.000000');
INSERT INTO `surplus_value_partitions` (`id`,`composition_id`,`total_surplus_value_lm`,`profit_enterprise_lm`,`interest_lm`,`ground_rent_lm`,`residual_surplus_lm`,`created_at`) VALUES
    ('5eed000000000000004902', '5eed000000000000004900', 1500, 900, 300, 300, 0, '1894-01-01 00:00:00.000000');

-- +goose Down
DELETE FROM `surplus_value_partitions` WHERE `id` = '5eed000000000000004902';
DELETE FROM `revenue_limits` WHERE `id` = '5eed000000000000004901';
DELETE FROM `annual_value_compositions` WHERE `id` = '5eed000000000000004900';
