-- +goose Up
-- +goose StatementBegin
ALTER TABLE Users ADD INDEX (Email)
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE Users ADD INDEX (Password)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE Users DROP INDEX Email
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE Users DROP INDEX Password
-- +goose StatementEnd
