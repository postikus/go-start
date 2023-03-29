-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS `reg_user` (
    `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
    `name` VARCHAR(80) NOT NULL DEFAULT ''
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS `user`;
-- +goose StatementEnd
