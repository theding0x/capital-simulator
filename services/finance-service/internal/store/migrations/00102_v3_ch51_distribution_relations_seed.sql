-- +goose Up
-- Ch. 51 — three production-relation bases (wage-labour, capital, landed property),
-- three distribution relations (wages→v, profit→s, rent→s) resting on them, one
-- historical-transience record for wages, one productive-force contradiction (19th-c crises).
INSERT INTO `production_relation_bases` (`id`,`name`,`characteristic_feature`,`is_historically_specific`,`created_at`) VALUES
    ('5eed000000000000005100', 'wage-labour', 'labour-power sold as a commodity', 1, '1894-01-01 00:00:00.000000'),
    ('5eed000000000000005106', 'capital', 'means of production confront workers as an alien power', 1, '1894-01-01 00:00:00.000000'),
    ('5eed000000000000005107', 'landed_property', 'monopoly of the earth', 1, '1894-01-01 00:00:00.000000');
INSERT INTO `distribution_relations` (`id`,`form`,`production_basis_id`,`apparent_character`,`real_character`,`corresponds_to_surplus_form`,`created_at`) VALUES
    ('5eed000000000000005101', 'wages', '5eed000000000000005100', 'a natural payment for work done', 'price of labour-power; historically specific to capitalism; disappears with wage-labour', 'v', '1894-01-01 00:00:00.000000'),
    ('5eed000000000000005102', 'profit', '5eed000000000000005106', 'a natural reward of capital', 'portion of surplus-value falling to the capitalist; regulates investment; historically specific', 's(profit)', '1894-01-01 00:00:00.000000'),
    ('5eed000000000000005103', 'rent', '5eed000000000000005107', 'a natural payment for the use of land', 'portion of surplus-value transferred to the landlord by the barrier to free capital investment; historically specific', 's(rent)', '1894-01-01 00:00:00.000000');
INSERT INTO `historical_transiences` (`id`,`distribution_form_id`,`mode_of_production`,`precondition`,`succeeding_form`,`created_at`) VALUES
    ('5eed000000000000005104', '5eed000000000000005101', 'capitalist', 'expropriation of the direct producers; labour-power becomes a commodity', 'distribution by need once wage-labour is abolished', '1894-01-01 00:00:00.000000');
INSERT INTO `productive_force_contradictions` (`id`,`productive_force`,`social_form`,`contradiction_nature`,`crisis_signal`,`created_at`) VALUES
    ('5eed000000000000005105', 'steam-powered socialised production, credit system, joint-stock companies', 'private property, anarchy of the market, wage-labour', 'socialised production constrained by private appropriation', 'commercial and industrial crises (1825, 1847, 1857)', '1894-01-01 00:00:00.000000');

-- +goose Down
DELETE FROM `productive_force_contradictions` WHERE `id` = '5eed000000000000005105';
DELETE FROM `historical_transiences` WHERE `id` = '5eed000000000000005104';
DELETE FROM `distribution_relations` WHERE `id` IN ('5eed000000000000005101','5eed000000000000005102','5eed000000000000005103');
DELETE FROM `production_relation_bases` WHERE `id` IN ('5eed000000000000005100','5eed000000000000005106','5eed000000000000005107');
