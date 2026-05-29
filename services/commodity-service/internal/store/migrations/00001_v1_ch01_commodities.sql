-- +goose Up
CREATE TABLE IF NOT EXISTS commodities (
    id               VARCHAR(24)  NOT NULL PRIMARY KEY,
    name             VARCHAR(255) NOT NULL,
    use_value_desc   TEXT         NOT NULL,
    use_value_unit   VARCHAR(100) NOT NULL,
    cl_kind          VARCHAR(100) NOT NULL DEFAULT '',
    cl_description   TEXT         NOT NULL,
    snlt_per_unit    BIGINT       NOT NULL DEFAULT 0,
    created_at       DATETIME(6)  NOT NULL,
    updated_at       DATETIME(6)  NOT NULL,
    UNIQUE KEY uq_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- +goose Down
DROP TABLE IF EXISTS commodities;
