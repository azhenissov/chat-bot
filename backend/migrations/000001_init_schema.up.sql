-- Включаем расширение для генерации UUID, если его нет
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Таблица чатов
CREATE TABLE chats (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Таблица сообщений
CREATE TABLE messages (
    id BIGSERIAL PRIMARY KEY,
    chat_id UUID NOT NULL REFERENCES chats(id) ON DELETE CASCADE,
    role VARCHAR(50) NOT NULL, -- 'user' или 'model'
    content TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Создаем индекс, так как мы будем часто искать сообщения по chat_id
CREATE INDEX idx_messages_chat_id ON messages(chat_id);