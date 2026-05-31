-- +goose Up
-- Seed: Ch. 33 — The Medium of Circulation in the Credit System
-- London Clearing House 1840s: bills of £150 000 000 settled with £1 000 000 in money.
INSERT INTO clearing_house_settlements
    (id, name, description, total_claims, money_used, net_balance, created_at)
VALUES
    ('5eed000000000000003300',
     'London Clearing House 1840s',
     'Bills of exchange circulate among bankers and merchants; mutual claims cancel in the clearing house; only the net balance requires actual coin. Marx cites this as proof that money functions merely as means of payment and measurement of value — the bulk of transactions require no actual currency.',
     15000000000,
     1000000000,
     14000000000,
     '1845-01-01 00:00:00.000000'),
    ('5eed000000000000003301',
     'Scotland velocity-9 note',
     'A £1 banknote issued on Monday by a Scottish country bank passes through nine hands — grain merchant, maltster, brewer, innkeeper, etc. — and returns to its issuing bank on Saturday. The same note performs nine payments; the effective supply is 9× the actual note stock.',
     900000,
     100000,
     800000,
     '1845-01-02 00:00:00.000000');

-- +goose Down
DELETE FROM clearing_house_settlements
WHERE id IN ('5eed000000000000003300', '5eed000000000000003301');
