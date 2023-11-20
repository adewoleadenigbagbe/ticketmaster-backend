-- +goose Up
-- +goose StatementBegin
ALTER TABLE Addresses
ADD IsCurrent TINYINT NULL
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE Addresses
DROP COLUMN IsCurrent
-- +goose StatementEnd
