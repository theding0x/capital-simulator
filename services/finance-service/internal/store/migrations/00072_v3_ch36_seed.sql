-- +goose Up
INSERT INTO usurers_capitals
    (id, name, stage, wage_labour, borrowers_json, interest_rate_bp, subordinated_to, description, created_at)
VALUES
    (
        '5eed000000000000003600',
        'Ancient Roman Usury',
        1,
        0,
        '["plebeians","small producers","nobles"]',
        4800,
        '',
        'Usury in ancient Rome operated at rates of 48% per annum or more, ruining plebeians and small producers alike. The creditor, a patrician or a wealthy equestrian, held the debtor''s person as security; default could mean debt bondage. No wage labour existed: the surplus was extracted by slavery and tribute, not the purchase of labour-power.',
        '0100-01-01 00:00:00.000000'
    ),
    (
        '5eed000000000000003601',
        'Medieval Church Usury',
        2,
        0,
        '["peasants","petty producers","feudal lords"]',
        2400,
        '',
        'Medieval usury operated beneath the formal prohibition of the Church. Montes pietatis, Jewish moneylenders, and the Italian merchant-bankers all charged rates of 24% or more. The borrowers were peasants burdened by bad harvests, petty producers needing tools, and feudal lords financing wars. The surplus was feudal ground-rent, not surplus-value; wage labour was embryonic at best.',
        '1300-01-01 00:00:00.000000'
    ),
    (
        '5eed000000000000003602',
        'Modern Subordinated Usury',
        4,
        1,
        '["industrial capitalists","commercial capitalists"]',
        800,
        'industrial_capital',
        'With the consolidation of the capitalist mode of production, usury capital was subordinated under industrial capital through the banking system. The rate of interest is now disciplined by the average rate of profit: usurers can no longer bleed the productive capitalist dry, for competition among money-lenders compresses interest to the fraction of profit that the industrial capitalist can afford to surrender. Wage labour is universal; surplus-value is the source of all interest.',
        '1870-01-01 00:00:00.000000'
    );

-- +goose Down
DELETE FROM usurers_capitals
WHERE id IN (
    '5eed000000000000003600',
    '5eed000000000000003601',
    '5eed000000000000003602'
);
