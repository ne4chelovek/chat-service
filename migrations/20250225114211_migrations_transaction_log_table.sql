-- +goose Up
-- +goose StatementBegin
CREATE TABLE transaction_log (
    id BIGSERIAL PRIMARY KEY,
    timestamp timestamp NOT NULL DEFAULT now(),
    log text NOT NULL
);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS transaction_log;
-- +goose StatementEnd
