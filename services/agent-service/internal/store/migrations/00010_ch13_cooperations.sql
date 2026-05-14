-- +goose Up
CREATE TABLE IF NOT EXISTS cooperations (
    id            VARCHAR(24)  NOT NULL PRIMARY KEY,
    name          VARCHAR(255) NOT NULL DEFAULT '',
    capitalist_id VARCHAR(24)  NOT NULL,
    members_json  JSON         NOT NULL,
    created_at    DATETIME(6)  NOT NULL,
    INDEX idx_cooperations_capitalist (capitalist_id)
);

-- +goose Down
DROP TABLE IF EXISTS cooperations;
