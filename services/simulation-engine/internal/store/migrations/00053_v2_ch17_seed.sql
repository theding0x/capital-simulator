-- +goose Up
-- Vol. II Ch. 17 — Seed: 100-spinning-mill aggregate fixture.
-- Seed IDs use prefix 5eed0000000000001701 (Ch. 17).

-- SurplusCirculation: 100-spinning-mill aggregate for period 1871.
-- Total surplus ≈ £208,333 = 50,000,000 pence.
INSERT INTO surplus_circulations
    (id, period, total_surplus_pence, realised_surplus_pence, unrealised_surplus_pence, created_at)
VALUES
    ('5eed0000000000001701', '1871', 50000000, 50000000, 0, NOW(6));

-- RealisationSourceMix: 80% consumption, 19% capitalised, 1% gold.
-- foreign_trade = 0 (Vol. II abstraction); credit = 0 for canonical fixture.
INSERT INTO realisation_source_entries
    (surplus_circulation_id, source, pence, deferred_realisation_pence)
VALUES
    ('5eed0000000000001701', 'capitalist_consumption', 40000000, 0),
    ('5eed0000000000001701', 'capitalised_surplus',     9500000, 0),
    ('5eed0000000000001701', 'gold_production',          500000, 0);

-- SocialCapitalAggregate snapshot for 1871.
-- department_i/ii_share_bp = 0 (placeholders until Ch. 20).
INSERT INTO social_capital_aggregates
    (period, total_advanced_pence, total_annual_output_pence,
     total_annual_surplus_pence, department_i_share_bp, department_ii_share_bp)
VALUES
    ('1871', 500000000, 550000000, 50000000, 0, 0);

-- Canonical RealisationPuzzle.
INSERT INTO realisation_puzzles
    (surplus_circulation_id, puzzle_statement, resolved_in_chapter)
VALUES
    ('5eed0000000000001701',
     'Where does the money come from to realise the aggregate surplus-value of all capitalists? Each individual capitalist can identify buyers for their own output, but at the social level the total money required to purchase C'' seems to exceed what was advanced as M. The puzzle dissolves only when we treat total social capital as a whole — formally resolved in the reproduction schemes of Chs. 20-21.',
     20);

-- +goose Down
DELETE FROM realisation_puzzles        WHERE surplus_circulation_id = '5eed0000000000001701';
DELETE FROM social_capital_aggregates  WHERE period = '1871';
DELETE FROM realisation_source_entries WHERE surplus_circulation_id = '5eed0000000000001701';
DELETE FROM surplus_circulations       WHERE id = '5eed0000000000001701';
