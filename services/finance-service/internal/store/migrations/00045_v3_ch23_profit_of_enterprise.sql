-- +goose Up
-- Vol. III Ch. 23 — Interest and Profit of Enterprise
-- profit_divisions stores ProfitDivision persisted entities.
-- TotalProfitBP > 0; InterestBP >= 0; InterestBP < TotalProfitBP.
-- ProfitOfEnterpriseBP = TotalProfitBP - InterestBP (partition invariant).
-- ownership_form: 1=OwnerOperator, 2=Separated.
-- mystification_index bp [0,10000].
CREATE TABLE `profit_divisions` (
    `id`                       VARCHAR(24)  NOT NULL,
    `total_profit_bp`          BIGINT       NOT NULL,
    `interest_bp`              BIGINT       NOT NULL,
    `profit_of_enterprise_bp`  BIGINT       NOT NULL,
    `ownership_form`           INT          NOT NULL,
    `wages_labour_minutes`     BIGINT       NOT NULL,
    `mystification_index`      BIGINT       NOT NULL,
    `created_at`               DATETIME(6)  NOT NULL,
    PRIMARY KEY (`id`),
    INDEX `idx_pd_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS `profit_divisions`;
