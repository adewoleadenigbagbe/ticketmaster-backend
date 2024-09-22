-- +goose Up
-- +goose StatementBegin
ALTER TABLE ShowSeats
MODIFY COLUMN BookingId CHAR(36) NULL
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE ShowSeats
MODIFY COLUMN BookingId CHAR(36) NOT NULL
-- +goose StatementEnd