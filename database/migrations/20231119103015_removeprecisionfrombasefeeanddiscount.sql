-- +goose Up
-- +goose StatementBegin
ALTER TABLE CinemaRates
MODIFY COLUMN BaseFee FLOAT NOT NULL
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE CinemaRates
MODIFY COLUMN Discount FLOAT
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
ALTER TABLE CinemaRates
MODIFY COLUMN BaseFee FLOAT(5,4) NOT NULL
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE CinemaRates
MODIFY COLUMN Discount FLOAT(4,3)
-- +goose StatementEnd

