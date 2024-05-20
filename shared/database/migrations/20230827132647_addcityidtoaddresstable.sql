-- +goose Up
-- +goose StatementBegin
ALTER TABLE Addresses
ADD CityId CHAR(36) NOT NULL,
ADD INDEX (CityId)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE Addresses
DROP COLUMN CityId,
DROP INDEX CityId
-- +goose StatementEnd
