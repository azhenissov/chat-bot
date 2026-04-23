package repository

import (
	"context"
	"encoding/json"
	"time"
	"chat-core/internal/domain"

	"github.com/redis/go-redis/v9"
)

type RedisRepo struct {
	client *redis.Client
}

func NewRedisRepo(client *redis.Client) *RedisRepo {
	return &RedisRepo{client: client}
}

// SaveContext сохраняет массив сообщений в Redis с временем жизни (TTL) 24 часа
func (r *RedisRepo) SaveContext(ctx context.Context, chatID string, contextMsgs []domain.Message) error {
	// Превращаем массив структур в JSON байты
	data, err := json.Marshal(contextMsgs)
	if err != nil {
		return err
	}

	// Ключ будет выглядеть как "chat_context:123e4567-e89b..."
	key := "chat_context:" + chatID
	return r.client.Set(ctx, key, data, 24*time.Hour).Err()
}

// GetContext достает последние сообщения для отправки в Python
func (r *RedisRepo) GetContext(ctx context.Context, chatID string) ([]domain.Message, error) {
	key := "chat_context:" + chatID
	data, err := r.client.Get(ctx, key).Bytes()
	
	if err == redis.Nil {
		// Ключ не найден (чат новый или кэш протух) - это не ошибка, просто возвращаем пустой массив
		return []domain.Message{}, nil
	} else if err != nil {
		return nil, err // Реальная ошибка сети или Redis
	}

	var messages []domain.Message
	if err := json.Unmarshal(data, &messages); err != nil {
		return nil, err
	}

	return messages, nil
}