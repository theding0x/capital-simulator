-- +goose Up
-- Vol. III Ch. 22 — Division of Profit; Rate of Interest; Natural Rate of Interest
-- rate_of_interests stores RateOfInterest persisted entities.
-- RateBP >= 0; RateBP < AverageProfitRateBP (interest is a portion of profit, strictly less).
-- CyclePhase ∈ {1=Inactivity, 2=Revival, 3=Prosperity, 4=Overproduction, 5=Crisis}.
CREATE TABLE `rate_of_interests` (
    `id`                      VARCHAR(24)  NOT NULL,
    `rate_bp`                 BIGINT       NOT NULL,
    `average_profit_rate_bp`  BIGINT       NOT NULL,
    `cycle_phase`             INT          NOT NULL,
    `period`                  VARCHAR(255) NOT NULL,
    `created_at`              DATETIME(6)  NOT NULL,
    PRIMARY KEY (`id`),
    INDEX `idx_roi_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS `rate_of_interests`;
