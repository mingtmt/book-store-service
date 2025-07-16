-- +goose Up
-- +goose StatementBegin
CREATE TABLE auths (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    username TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT now()
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS auths;
-- +goose StatementEnd