-- +goose Up
-- Atlas: per-capital turnover number drives the orbit dot's lap rate (Vol. II Ch.7).
ALTER TABLE industrial_capitals ADD COLUMN turnover_number INT NOT NULL DEFAULT 1;
UPDATE industrial_capitals SET turnover_number = 5 WHERE id = '5eed000000000000000401'; -- Spinning Mill (5×/yr)
UPDATE industrial_capitals SET turnover_number = 8 WHERE id = '5eed000000000000009001'; -- Tannery (light, fast)
UPDATE industrial_capitals SET turnover_number = 2 WHERE id = '5eed000000000000009011'; -- Steelworks (heavy fixed capital, slow)
UPDATE industrial_capitals SET turnover_number = 4 WHERE id = '5eed000000000000009021'; -- Textile

-- +goose Down
ALTER TABLE industrial_capitals DROP COLUMN turnover_number;
