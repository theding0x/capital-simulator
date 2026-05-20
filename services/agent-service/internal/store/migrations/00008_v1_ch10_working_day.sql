-- +goose Up
CREATE TABLE IF NOT EXISTS working_days (
    id                       VARCHAR(24)  NOT NULL PRIMARY KEY,
    necessary_labour_minutes BIGINT       NOT NULL,
    surplus_labour_minutes   BIGINT       NOT NULL,
    created_at               DATETIME(6)  NOT NULL
);

CREATE TABLE IF NOT EXISTS relay_schedules (
    id              VARCHAR(24)  NOT NULL PRIMARY KEY,
    shift_kind_0    VARCHAR(10)  NOT NULL,
    nl_minutes_0    BIGINT       NOT NULL,
    sl_minutes_0    BIGINT       NOT NULL,
    worker_ids_0    JSON         NOT NULL,
    shift_kind_1    VARCHAR(10)  NOT NULL,
    nl_minutes_1    BIGINT       NOT NULL,
    sl_minutes_1    BIGINT       NOT NULL,
    worker_ids_1    JSON         NOT NULL,
    created_at      DATETIME(6)  NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS relay_schedules;
DROP TABLE IF EXISTS working_days;
