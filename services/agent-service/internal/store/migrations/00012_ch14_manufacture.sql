-- +goose Up
-- Ch. 14 — Division of Labour and Manufacture.
--
-- A Manufacture is a Cooperation (Ch. 13) reorganised around a fixed
-- division of detail labour: it has a form (heterogeneous / serial), a
-- two-fold origin (combination / splitting), and a list of DetailRoles
-- each carrying a skill level, an hourly output rate, a headcount, and a
-- specialised tool.

CREATE TABLE IF NOT EXISTS manufactures (
    id                             VARCHAR(24)  NOT NULL PRIMARY KEY,
    name                           VARCHAR(255) NOT NULL DEFAULT '',
    capitalist_id                  VARCHAR(24)  NOT NULL,
    form                           VARCHAR(32)  NOT NULL,
    origin                         VARCHAR(32)  NOT NULL,
    individual_working_day_minutes BIGINT       NOT NULL DEFAULT 0,
    created_at                     DATETIME(6)  NOT NULL,
    INDEX idx_manufactures_capitalist (capitalist_id),
    INDEX idx_manufactures_form       (form)
);

CREATE TABLE IF NOT EXISTS detail_roles (
    manufacture_id       VARCHAR(24)  NOT NULL,
    ordinal              INT          NOT NULL,
    role_name            VARCHAR(255) NOT NULL,
    skill_level          VARCHAR(32)  NOT NULL,
    output_rate_per_hour BIGINT       NOT NULL DEFAULT 0,
    head_count           INT          NOT NULL DEFAULT 0,
    tool_name            VARCHAR(255) NOT NULL DEFAULT '',
    PRIMARY KEY (manufacture_id, ordinal),
    INDEX idx_detail_roles_manufacture (manufacture_id)
);

-- +goose Down
DROP TABLE IF EXISTS detail_roles;
DROP TABLE IF EXISTS manufactures;
