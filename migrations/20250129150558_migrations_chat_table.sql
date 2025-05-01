-- +goose Up
-- +goose StatementBegin

CREATE TABLE chats (
    chat_id BIGSERIAL PRIMARY KEY,     -- уникальный идентификатор чата
    usernames text[] NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP  -- время создания чата
);

CREATE TABLE messages (
    chat_id BIGINT NOT NULL REFERENCES chats(chat_id) ON DELETE CASCADE,  -- ссылка на чат
    from_user VARCHAR(100) NOT NULL,  -- отправитель (ссылка на пользователя из сервиса авторизации)
    text TEXT NOT NULL,  -- текст сообщения
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,  -- время отправки
    status VARCHAR(20) CHECK (status IN ('SENT', 'DELIVERED', 'READ')) DEFAULT 'SENT'  -- статус сообщения
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS chats;
DROP TABLE IF EXISTS messages;

-- +goose StatementEnd
