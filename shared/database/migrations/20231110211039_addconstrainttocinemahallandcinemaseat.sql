-- +goose Up
-- +goose StatementBegin
ALTER TABLE CinemaHalls 
ADD CONSTRAINT fk_cinemahall_cinema
FOREIGN KEY (CinemaId) REFERENCES Cinemas(Id) ON DELETE CASCADE ON UPDATE CASCADE;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE CinemaSeats 
ADD CONSTRAINT fk_cinemaseat_cinemahall  
FOREIGN KEY (CinemaHallId) REFERENCES CinemaHalls(Id) ON DELETE CASCADE ON UPDATE CASCADE;
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
ALTER TABLE CinemaHalls
DROP FOREIGN KEY fk_cinemahall_cinema
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE CinemaSeats
DROP FOREIGN KEY fk_cinemaseat_cinemahall
-- +goose StatementEnd

