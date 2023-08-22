-- +goose Up
-- +goose StatementBegin
CREATE TABLE Bookings (
 Id CHAR(36) PRIMARY KEY,
 NumberOfSeats INT NOT NULL,
 BookDateTime DATETIME NOT NULL,
 Status int NOT NULL,
 UserId CHAR(36) NOT NULL,
 ShowId CHAR(36) NOT NULL,
 IsDeprecated TINYINT,
 CreatedOn NOT NULL DATETIME DEFAULT,
 ModifiedOn NOT NULL DATETIME ON UPDATE,
 INDEX (UserId),
 INDEX (ShowId),
 INDEX (CreatedOn)
);

CREATE TABLE Cinemas (
 Id CHAR(36) PRIMARY KEY,
 Name VARCHAR(255) NOT NULL,
 TotalCinemalHalls INT NOT NULL,
 CityId CHAR(36) NOT NULL,
 IsDeprecated TINYINT,
 CreatedOn NOT NULL DATETIME DEFAULT,
 ModifiedOn NOT NULL DATETIME ON UPDATE,
 INDEX (Name),
 INDEX (CityId),
 INDEX (CreatedOn)
);

CREATE TABLE CinemaHalls (
 Id CHAR(36) PRIMARY KEY,
 Name VARCHAR(255) NOT NULL,
 TotalSeats INT NOT NULL,
 CinemaId CHAR(36) NOT NULL,
 IsDeprecated TINYINT,
 CreatedOn NOT NULL DATETIME DEFAULT,
 ModifiedOn NOT NULL DATETIME ON UPDATE,
 INDEX (Name),
 INDEX (CinemaId),
 INDEX (CreatedOn)
);

CREATE TABLE CinemaSeats (
 Id CHAR(36) PRIMARY KEY,
 SeatNumber INT NOT NULL,
 Type INT NOT NULL,
 CinemaHallId CHAR(36) NOT NULL,
 IsDeprecated TINYINT,
 CreatedOn NOT NULL DATETIME DEFAULT,
 ModifiedOn NOT NULL DATETIME ON UPDATE,
 INDEX (CinemaHallId),
 INDEX (CreatedOn)
);

CREATE TABLE Cities (
 Id CHAR(36) PRIMARY KEY,
 Name VARCHAR(255) NOT NULL,
 State VARCHAR(255) NOT NULL,
 ZipCode VARCHAR(255) NOT NULL,
 CinemaHallId CHAR(36) NOT NULL,
 IsDeprecated TINYINT,
 CreatedOn NOT NULL DATETIME DEFAULT,
 ModifiedOn NOT NULL DATETIME ON UPDATE,
 INDEX (Name),
 INDEX (CreatedOn)
);

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
 IsDeprecated TINYINT,
 CreatedOn NOT NULL DATETIME DEFAULT,
 ModifiedOn NOT NULL DATETIME ON UPDATE,
 INDEX (Name),
 INDEX (CreatedOn)
);

CREATE TABLE Payments (
 Id CHAR(36) PRIMARY KEY,
 PaymentDate DATETIME NOT NULL,
 DiscountCouponId CHAR(36),
 RemoteTransactionId CHAR(36),
 PaymentMethod INT NOT NULL,
 Amount DECIMAL(19, 4) NOT NULL,
 BookingId CHAR(36) NOT NULL,
 IsDeprecated TINYINT,
 CreatedOn NOT NULL DATETIME DEFAULT,
 ModifiedOn NOT NULL DATETIME ON UPDATE,
 INDEX (PaymentDate),
 INDEX (CreatedOn)
);

CREATE TABLE Shows (
 Id CHAR(36) PRIMARY KEY,
 Date DATETIME NOT NULL,
 StartTime INT(64) NULL,
 EndTime INT(64) NULL,
 PaymentMethod INT NOT NULL,
 Amount DECIMAL(19, 4) NOT NULL,
 CinemaHallId CHAR(36) NOT NULL,
 MovieId CHAR(36) NOT NULL,
 IsDeprecated TINYINT,
 IsCancelled TINYINT,
 CancellationReason VARCHAR(255)
 CreatedOn NOT NULL DATETIME DEFAULT,
 ModifiedOn NOT NULL DATETIME ON UPDATE,
 INDEX (Date),
 INDEX (CinemaHallId)
 INDEX (MovieId),
 INDEX (CreatedOn)
);

CREATE TABLE ShowSeats (
 Id CHAR(36) PRIMARY KEY,
 Status INT NOT NULL,
 Price DECIMAL(19, 4) NOT NULL,
 CinemaSeatId CHAR(36) NOT NULL,
 ShowId CHAR(36) NOT NULL,
 BookingId CHAR(36) NOT NULL,
 IsDeprecated TINYINT,
 CreatedOn NOT NULL DATETIME DEFAULT,
 ModifiedOn NOT NULL DATETIME ON UPDATE,
 INDEX (CinemaSeatId)
 INDEX (ShowId),
 INDEX (BookingId),
 INDEX (CreatedOn)
);

CREATE TABLE Users (
 Id CHAR(36) PRIMARY KEY,
 FirstName VARCHAR(255) NOT NULL,
 LastName VARCHAR(255) NOT NULL,
 Email VARCHAR(50) NOT NULL,
 PhoneNumber VARCHAR(20) NOT NULL,
 Password VARCHAR(255) NOT NULL,
 Price DECIMAL(19, 4) NOT NULL,
 CreatedOn NOT NULL DATETIME DEFAULT,
 ModifiedOn NOT NULL DATETIME ON UPDATE,
 INDEX (CreatedOn)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE Bookings;
DROP TABLE Cinemas;
DROP TABLE CinemaHalls;
DROP TABLE CinemaSeats;
DROP TABLE Cities;
DROP TABLE Movies;
DROP TABLE Payments;
DROP TABLE Shows;
DROP TABLE ShowSeats;
DROP TABLE Users;
-- +goose StatementEnd
