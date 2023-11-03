-- +goose Up

-- +goose StatementBegin
CREATE TABLE Bookings (
 Id CHAR(36) PRIMARY KEY,
 NumberOfSeats INT NOT NULL,
 BookDateTime DATETIME NOT NULL,
 Status int NOT NULL,
 UserId CHAR(36) NOT NULL,
 ShowId CHAR(36) NOT NULL,
 IsDeprecated TINYINT NOT NULL,
 CreatedOn DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
 ModifiedOn DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
 INDEX (UserId),
 INDEX (ShowId),
 INDEX (CreatedOn)
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE Cinemas (
 Id CHAR(36) PRIMARY KEY,
 Name VARCHAR(255) NOT NULL,
 TotalCinemalHalls INT NOT NULL,
 CityId CHAR(36) NOT NULL,
 IsDeprecated TINYINT NOT NULL,
 CreatedOn DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
 ModifiedOn DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
 INDEX (Name),
 INDEX (CityId),
 INDEX (CreatedOn)
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE CinemaHalls (
 Id CHAR(36) PRIMARY KEY,
 Name VARCHAR(255) NOT NULL,
 TotalSeat INT NOT NULL,
 CinemaId CHAR(36) NOT NULL,
 IsDeprecated TINYINT NOT NULL,
 CreatedOn DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
 ModifiedOn DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
 INDEX (Name),
 INDEX (CinemaId),
 INDEX (CreatedOn)
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE CinemaSeats (
 Id CHAR(36) PRIMARY KEY,
 SeatNumber INT NOT NULL,
 Type INT NOT NULL,
 CinemaHallId CHAR(36) NOT NULL,
 IsDeprecated TINYINT NOT NULL,
 CreatedOn DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
 ModifiedOn DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
 INDEX (CinemaHallId),
 INDEX (CreatedOn)
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE Cities (
 Id CHAR(36) PRIMARY KEY,
 Name VARCHAR(255) NOT NULL,
 State VARCHAR(255) NOT NULL,
 ZipCode VARCHAR(255) NOT NULL,
 IsDeprecated TINYINT NOT NULL,
 CreatedOn DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
 ModifiedOn DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
 INDEX (Name)
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE Movies (
 Id CHAR(36) PRIMARY KEY,
 Title MEDIUMTEXT NOT NULL,
 Description LONGTEXT,
 Language CHAR(10) NOT NULL,
 ReleaseDate DATETIME NOT NULL,
 Duration INT,
 Genre INT NOT NULL,
 Popularity FLOAT NOT NULL,
 VoteCount INT NOT NULL,
 IsDeprecated TINYINT NOT NULL,
 CreatedOn DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
 ModifiedOn DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
 INDEX (CreatedOn)
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE Payments (
 Id CHAR(36) PRIMARY KEY,
 Amount DECIMAL(19, 4) NOT NULL,
 PaymentDate DATETIME NOT NULL,
 DiscountCouponId CHAR(36),
 RemoteTransactionId CHAR(36),
 PaymentMethod INT NOT NULL,
 BookingId CHAR(36) NOT NULL,
 IsDeprecated TINYINT NOT NULL,
 CreatedOn DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
 ModifiedOn DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
 INDEX (PaymentDate),
 INDEX (BookingId),
 INDEX (CreatedOn)
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE Shows (
 Id CHAR(36) PRIMARY KEY,
 Date DATETIME NOT NULL,
 StartTime INT(64) NOT NULL,
 EndTime INT(64) NOT NULL,
 CinemaHallId CHAR(36) NOT NULL,
 MovieId CHAR(36) NOT NULL,
 IsDeprecated TINYINT NOT NULL,
 IsCancelled TINYINT,
 CancellationReason VARCHAR(255),
 CreatedOn DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
 ModifiedOn DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
 INDEX (Date),
 INDEX (CinemaHallId),
 INDEX (MovieId),
 INDEX (CreatedOn)
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE ShowSeats (
 Id CHAR(36) PRIMARY KEY,
 Status INT NOT NULL,
 Price DECIMAL(19, 4) NOT NULL,
 CinemaSeatId CHAR(36) NOT NULL,
 ShowId CHAR(36) NOT NULL,
 BookingId CHAR(36) NOT NULL,
 IsDeprecated TINYINT NOT NULL,
 CreatedOn DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
 ModifiedOn DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
 INDEX (CinemaSeatId),
 INDEX (ShowId),
 INDEX (BookingId),
 INDEX (CreatedOn)
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE Users (
 Id CHAR(36) PRIMARY KEY,
 FirstName VARCHAR(255) NOT NULL,
 LastName VARCHAR(255) NOT NULL,
 Email VARCHAR(50) NOT NULL,
 PhoneNumber VARCHAR(20),
 Password VARCHAR(255) NOT NULL,
 IsDeprecated TINYINT NOT NULL,
 CreatedOn DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
 ModifiedOn DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
 INDEX (CreatedOn)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE Bookings;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE Cinemas;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE CinemaHalls;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE CinemaSeats;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE Cities;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE Movies;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE Payments;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE Shows;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE ShowSeats;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE Users;
-- +goose StatementEnd