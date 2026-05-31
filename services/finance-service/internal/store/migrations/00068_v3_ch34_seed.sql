-- +goose Up
INSERT INTO note_issue_constraints
    (id, name, regime, max_notes, gold_backing, unbacked_limit, notes_beyond_max, description, created_at)
VALUES
    ('5eed000000000000003400', 'Bank Act 1844 — founding regime', 2, 2800000000, 1400000000, 1400000000, 0,
     'Bank Charter Act 1844 enacted: Issue Department limited to £14m fiduciary (unbacked) notes; additional notes require 100% gold backing. Total statutory maximum £28m. Currency School doctrine institutionalised.',
     '1844-07-19 00:00:00.000000'),
    ('5eed000000000000003401', 'Crisis of 1847 — Act in force', 2, 2800000000, 1400000000, 1400000000, 0,
     'Commercial crisis of 1847: Act remains formally in force but Treasury prepares suspension letter. Gold drain accompanies high interest rates, not high commodity prices — empirically challenging the Currency Principle.',
     '1847-10-01 00:00:00.000000'),
    ('5eed000000000000003402', 'Crisis of 1857 — Act suspended', 3, 2800000000, 1400000000, 1400000000, 48883000,
     'Commercial crisis of 1857: Bank Act 1844 suspended by Treasury letter on 12 November. Bank authorised to issue notes beyond the £28m statutory maximum; £48,883 in excess notes actually issued before crisis abated without the suspension being exercised in full.',
     '1857-11-12 00:00:00.000000');

-- +goose Down
DELETE FROM note_issue_constraints
WHERE id IN (
    '5eed000000000000003400',
    '5eed000000000000003401',
    '5eed000000000000003402'
);
