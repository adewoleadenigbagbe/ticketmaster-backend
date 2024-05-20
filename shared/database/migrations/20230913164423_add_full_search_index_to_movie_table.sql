-- +goose Up
-- +goose StatementBegin
ALTER TABLE Movies 
ADD FULLTEXT INDEX Idx_title_Description(title, Description);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE Movies
DROP INDEX Idx_title_Description;
-- +goose StatementEnd
