-- +goose Up
-- Ch. 52 — the three great classes (wage-labourers/capitalists/landowners), their income
-- sources (wages→v component 0; profit/rent→s components 900/600), their developmental
-- tendencies, and the officials ambiguity whose resolution the manuscript never reaches.
INSERT INTO `social_classes` (`id`,`name`,`relation_to_means_of_production`,`primary_revenue_form`,`historical_precondition`,`created_at`) VALUES
    ('5eed000000000000005200', 'wage-labourers', 'separated from the means of production; sell labour-power as a commodity', 'wages', 'primitive accumulation: expropriation of the peasants from the land', '1894-01-01 00:00:00.000000'),
    ('5eed000000000000005201', 'capitalists', 'own the means of production as capital confronting living labour', 'profit', 'accumulation of money-capital; concentration of the means of production', '1894-01-01 00:00:00.000000'),
    ('5eed000000000000005202', 'landowners', 'monopoly ownership of the earth; collect tribute without producing', 'rent', 'private ownership of land historically established', '1894-01-01 00:00:00.000000');
INSERT INTO `class_income_sources` (`id`,`class_id`,`revenue_form`,`surplus_value_component_bp`,`created_at`) VALUES
    ('5eed000000000000005205', '5eed000000000000005200', 'wages', 0, '1894-01-01 00:00:00.000000'),
    ('5eed000000000000005206', '5eed000000000000005201', 'profit', 900, '1894-01-01 00:00:00.000000'),
    ('5eed000000000000005207', '5eed000000000000005202', 'rent', 600, '1894-01-01 00:00:00.000000');
INSERT INTO `class_tendencies` (`id`,`class_id`,`description`,`direction`,`created_at`) VALUES
    ('5eed000000000000005203', '5eed000000000000005201', 'fewer, larger capitals as the means of production concentrate', 'concentration', '1894-01-01 00:00:00.000000'),
    ('5eed000000000000005208', '5eed000000000000005200', 'more workers as expropriation proceeds', 'expansion', '1894-01-01 00:00:00.000000'),
    ('5eed000000000000005209', '5eed000000000000005202', 'landed property takes the capitalist form', 'subordination', '1894-01-01 00:00:00.000000');
INSERT INTO `class_ambiguities` (`id`,`intermediate_group`,`income_source`,`why_not_a_class`,`created_at`) VALUES
    ('5eed000000000000005204', 'officials', 'salary', 'manuscript_breaks_off', '1894-01-01 00:00:00.000000');

-- +goose Down
DELETE FROM `class_ambiguities` WHERE `id` = '5eed000000000000005204';
DELETE FROM `class_tendencies` WHERE `id` IN ('5eed000000000000005203','5eed000000000000005208','5eed000000000000005209');
DELETE FROM `class_income_sources` WHERE `id` IN ('5eed000000000000005205','5eed000000000000005206','5eed000000000000005207');
DELETE FROM `social_classes` WHERE `id` IN ('5eed000000000000005200','5eed000000000000005201','5eed000000000000005202');
