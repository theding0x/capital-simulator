-- +goose Up
CREATE TABLE note_issue_constraints (
    id              VARCHAR(24)  NOT NULL,
    name            VARCHAR(255) NOT NULL DEFAULT '',
    regime          INT          NOT NULL DEFAULT 1,
    max_notes       BIGINT       NOT NULL DEFAULT 0,
    gold_backing    BIGINT       NOT NULL DEFAULT 0,
    unbacked_limit  BIGINT       NOT NULL DEFAULT 0,
    notes_beyond_max BIGINT      NOT NULL DEFAULT 0,
    description     TEXT         NOT NULL,
    created_at      DATETIME(6)  NOT NULL,
    PRIMARY KEY (id),
    INDEX idx_note_issue_constraints_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS note_issue_constraints;
