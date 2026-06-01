-- +goose Up
-- Issue #220: record the labour-force reproduction check on each tick. Each
-- period's wage bill (the scheme's total variable capital) is asserted to cover
-- the worker pool's subsistence (worker_pool_size * subsistence_basket_pence).
-- worker_pool_size = 0 means the check was not requested for that tick.
-- tick_status is "reproduced" or "labour-force-shortfall".
ALTER TABLE reproduction_ticks
    ADD COLUMN worker_pool_size         BIGINT      NOT NULL DEFAULT 0,
    ADD COLUMN wage_bill_pence          BIGINT      NOT NULL DEFAULT 0,
    ADD COLUMN subsistence_basket_pence BIGINT      NOT NULL DEFAULT 0,
    ADD COLUMN subsistence_covered      TINYINT(1)  NOT NULL DEFAULT 1,
    ADD COLUMN tick_status              VARCHAR(32) NOT NULL DEFAULT 'reproduced';

-- +goose Down
ALTER TABLE reproduction_ticks
    DROP COLUMN tick_status,
    DROP COLUMN subsistence_covered,
    DROP COLUMN subsistence_basket_pence,
    DROP COLUMN wage_bill_pence,
    DROP COLUMN worker_pool_size;
