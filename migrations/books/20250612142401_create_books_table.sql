-- +goose Up
-- +goose StatementBegin
CREATE TABLE books (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    title TEXT NOT NULL,
    author TEXT NOT NULL,
    price NUMERIC NOT NULL,
    created_at TIMESTAMP DEFAULT now()
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE books;
-- +goose StatementEnd