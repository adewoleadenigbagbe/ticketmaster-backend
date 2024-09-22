-- +goose Up
-- +goose StatementBegin
CREATE TABLE CinemaRates (
 Id CHAR(36) PRIMARY KEY,
 CinemaId CHAR(36) NOT NULL,
 BaseFee FLOAT(5,4) NOT NULL,
 IsActive TINYINT NOT NULL,
 Discount FLOAT(4,3) NULL,
 IsSpecials TINYINT NULL,
 CreatedOn DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
 ModifiedOn DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
 INDEX (CinemaId)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE CinemaRates;
-- +goose StatementEnd
