-- +goose Up
-- +goose StatementBegin
ALTER TABLE Payments
ADD ProviderExtraInformation LONGTEXT NULL
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE Payments
DROP COLUMN ProviderExtraInformation
-- +goose StatementEnd
