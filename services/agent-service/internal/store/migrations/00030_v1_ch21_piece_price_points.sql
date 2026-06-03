-- +goose Up
-- Capital Vol. I, Ch. 21 (#219) — the dynamic piece-price time series.
--
-- As factory productivity (Ch. 15) rises, the daily wage is unchanged and the
-- socially-normal output grows, so the price per piece falls in inverse
-- proportion. The simulation-engine tick scheduler appends one row here each
-- time the productivity factor driving a piece-wage's worker changes.
--
-- Note: `sequence` is a MySQL reserved word, so the per-agent ordering column
-- is named sequence_num.
CREATE TABLE piece_price_points (
    id                  CHAR(24)    NOT NULL PRIMARY KEY,
    agent_id            CHAR(24)    NOT NULL,
    sequence_num        BIGINT      NOT NULL,
    productivity_factor DOUBLE      NOT NULL,
    price_pence         BIGINT      NOT NULL,
    normal_output       BIGINT      NOT NULL,
    occurred_at         DATETIME(6) NOT NULL,
    UNIQUE KEY uk_piece_price_points_agent_seq (agent_id, sequence_num),
    KEY idx_piece_price_points_agent (agent_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS piece_price_points;
