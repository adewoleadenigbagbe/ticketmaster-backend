-- +goose Up
-- +goose StatementBegin
CREATE TABLE Addresses (
 Id CHAR(36) PRIMARY KEY,
 AddressLine MEDIUMTEXT NOT NULL,
 AddressType INT NOT NULL,
 Coordinates POINT NOT NULL,
 EntityId CHAR(36) NOT NULL,
 IsDeprecated TINYINT NOT NULL,
 CreatedOn DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
 ModifiedOn DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
 INDEX (EntityId),
 INDEX (CreatedOn)
);
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE Cities
ADD Coordinates POINT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE Addresses;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE Cities
DROP COLUMN Coordinates;
-- +goose StatementEnd
