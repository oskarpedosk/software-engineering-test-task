-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN uuid UUID DEFAULT gen_random_uuid() UNIQUE NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN uuid;
-- +goose StatementEnd
