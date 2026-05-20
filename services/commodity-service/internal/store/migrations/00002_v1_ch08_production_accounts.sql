-- +goose Up
-- Ch. 8-9: production accounts record c, v, s for a single production run.
-- constant/variable/surplus are stored as BIGINT labour-minutes.
CREATE TABLE IF NOT EXISTS production_accounts (
    id         VARCHAR(24) NOT NULL PRIMARY KEY,
    constant   BIGINT      NOT NULL DEFAULT 0,
    variable   BIGINT      NOT NULL DEFAULT 0,
    surplus    BIGINT      NOT NULL DEFAULT 0,
    created_at DATETIME(6) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- +goose Down
DROP TABLE IF EXISTS production_accounts;
