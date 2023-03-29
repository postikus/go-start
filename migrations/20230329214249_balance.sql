-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS `balance` (
    `user_id` BIGINT PRIMARY KEY,
    `amount` DECIMAL NOT NULL DEFAULT 0,
    CONSTRAINT `fk_balance_reg_user_id` FOREIGN KEY (`user_id`) REFERENCES `reg_user`(`id`)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS `balance`;
-- +goose StatementEnd
