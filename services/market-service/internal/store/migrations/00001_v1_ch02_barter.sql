-- +goose Up
CREATE TABLE IF NOT EXISTS owners (
    id         VARCHAR(24)  NOT NULL PRIMARY KEY,
    name       VARCHAR(255) NOT NULL,
    created_at DATETIME(6)  NOT NULL,
    updated_at DATETIME(6)  NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS offers (
    id                 VARCHAR(24)  NOT NULL PRIMARY KEY,
    owner_id           VARCHAR(24)  NOT NULL,
    commodity_id       VARCHAR(24)  NOT NULL,
    quantity           DOUBLE       NOT NULL,
    seeks_kind         VARCHAR(20)  NOT NULL DEFAULT '',
    seeks_commodity_id VARCHAR(24)  NOT NULL DEFAULT '',
    created_at         DATETIME(6)  NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS exchanges (
    id                     VARCHAR(24) NOT NULL PRIMARY KEY,
    giver_id               VARCHAR(24) NOT NULL,
    receiver_id            VARCHAR(24) NOT NULL,
    giver_commodity_id     VARCHAR(24) NOT NULL,
    giver_qty              DOUBLE      NOT NULL,
    receiver_commodity_id  VARCHAR(24) NOT NULL,
    receiver_qty           DOUBLE      NOT NULL,
    realised_value         BIGINT      NOT NULL,
    created_at             DATETIME(6) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS exchanges;
DROP TABLE IF EXISTS offers;
DROP TABLE IF EXISTS owners;
