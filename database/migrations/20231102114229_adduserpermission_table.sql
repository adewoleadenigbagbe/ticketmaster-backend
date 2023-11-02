-- +goose Up
-- +goose StatementBegin
CREATE TABLE UserRoles (
 Id CHAR(36) PRIMARY KEY,
 Name VARCHAR(255) NOT NULL,
 Description MEDIUMTEXT NOT NULL,
 Role INT NOT NULL,
 CreatedOn DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
 ModifiedOn DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
 INDEX (CreatedOn)
);
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE Users 
ADD RoleId CHAR(36) NOT NULL,
ADD CONSTRAINT fk_user_userRole  
FOREIGN KEY (RoleId) REFERENCES UserRoles(Id) ON DELETE RESTRICT ON UPDATE RESTRICT;
-- +goose StatementEnd

-- +goose Down

-- +goose StatementBegin
ALTER TABLE Users
DROP COLUMN RoleId,
DROP FOREIGN KEY fk_user_userRole
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE UserRoles;
-- +goose StatementEnd
