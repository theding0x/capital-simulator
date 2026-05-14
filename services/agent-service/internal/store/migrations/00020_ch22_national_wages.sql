-- +goose Up
CREATE TABLE national_intensities (
    country_code CHAR(2) NOT NULL PRIMARY KEY,
    factor       DOUBLE  NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE day_wages (
    country_code        CHAR(2) NOT NULL PRIMARY KEY,
    nominal_pence       BIGINT  NOT NULL,
    working_day_minutes BIGINT  NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE spindle_ratios (
    country_code        CHAR(2) NOT NULL PRIMARY KEY,
    spindles_per_worker BIGINT  NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS spindle_ratios;
DROP TABLE IF EXISTS day_wages;
DROP TABLE IF EXISTS national_intensities;
