-- +goose Up
ALTER TABLE labour_workers
    ADD COLUMN labour_power_value_minutes BIGINT NOT NULL DEFAULT 0;

CREATE TABLE IF NOT EXISTS labour_processes (
    id                        VARCHAR(24)   NOT NULL PRIMARY KEY,
    worker_id                 VARCHAR(24)   NOT NULL,
    capitalist_id             VARCHAR(24)   NOT NULL,
    duration                  BIGINT        NOT NULL,
    necessary_labour_minutes  BIGINT        NOT NULL DEFAULT 0,
    means_json                JSON          NOT NULL,
    product_kind              VARCHAR(255)  NOT NULL DEFAULT '',
    product_quantity          BIGINT        NOT NULL DEFAULT 0,
    created_at                DATETIME(6)   NOT NULL,
    INDEX idx_lp_worker_id    (worker_id),
    INDEX idx_lp_capitalist   (capitalist_id)
);

-- +goose Down
DROP TABLE IF EXISTS labour_processes;
ALTER TABLE labour_workers DROP COLUMN labour_power_value_minutes;
