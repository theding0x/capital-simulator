-- +goose Up
-- Vol. III Ch. 30 seeds — Marx's money-capital vs real-capital examples.
-- Seed IDs start with 5eed000000000000003000..3001 (chapter 30).
--
-- Seed 3000: Crisis 1847 (phase=5) — loanable money piles up in the banks while
--            real capital (productive plant + saleable stocks) contracts; the two
--            accumulations diverge (corresponds=0). Commercial credit contracts,
--            but NOT because money is scarce: returns are blocked (commodities
--            will not sell) and mutual confidence has collapsed.
--
-- Seed 3001: Revival 1850s (phase=2) — loanable funds and productive investment
--            expand together (corresponds=1); commercial credit flows freely; the
--            reserve is ample. Money- and real-capital accumulation correspond.

INSERT INTO `real_capital_accumulations`
    (`id`, `name`, `phase`, `money_capital_growth_bp`, `real_capital_growth_bp`, `corresponds`, `correspondence_explanation`, `reserve_capital`, `credit_limit_description`, `commercial_credit_contracts`, `cause_is_money_scarcity`, `description`, `created_at`)
VALUES
    ('5eed000000000000003000',
     'Crisis 1847',
     5,
     800,
     -600,
     0,
     'loanable money piles up in banks while productive capital and saleable stocks contract',
     50000000,
     'Bank of England banking reserve — the hard floor in a crisis',
     1,
     0,
     'commercial credit contracts because returns are blocked and confidence is lost, not from a money shortage; the ultimate cause is restricted mass consumption against over-extended production',
     '2025-01-01 00:00:00.000000'),
    ('5eed000000000000003001',
     'Revival 1850s',
     2,
     400,
     450,
     1,
     'loanable funds and productive investment expand together',
     300000000,
     'ample reserve in the revival phase',
     0,
     0,
     'money- and real-capital accumulation correspond; commercial credit flows freely on the rhythm of reproduction',
     '2025-01-01 00:01:00.000000');

-- +goose Down
DELETE FROM `real_capital_accumulations`
WHERE `id` IN ('5eed000000000000003000', '5eed000000000000003001');
