// core-backend/internal/service/chat.go
package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"chat-core/internal/domain"
	"chat-core/internal/repository"
)

type PythonChatRequest struct {
	ChatID     string           `json:"chat_id"`
	NewMessage string           `json:"new_message"`
	Context    []domain.Message `json:"context"` 
}

type PythonChatResponse struct {
	Response string `json:"response"`
}

type ChatService struct {
	pgRepo  *repository.PostgresRepo
	rdRepo  *repository.RedisRepo
	aiURL   string // URL питоновского сервиса, например "http://localhost:8000/generate"
}

func NewChatService(pg *repository.PostgresRepo, rd *repository.RedisRepo, aiURL string) *ChatService {
	return &ChatService{
		pgRepo: pg,
		rdRepo: rd,
		aiURL:  aiURL,
	}
}


func (s *ChatService) ProcessMessage(ctx context.Context, chatID string, userText string) (string, error) {
	// 1. Сохраняем сообщение пользователя в Postgres
	userMsg := &domain.Message{
		ChatID:  chatID,
		Role:    "user",
		Content: userText,
	}
	if err := s.pgRepo.SaveMessage(ctx, userMsg); err != nil {
		return "", fmt.Errorf("ошибка сохранения сообщения юзера: %w", err)
	}

	// 2. Достаем историю (контекст) из Redis
	contextMsgs, err := s.rdRepo.GetContext(ctx, chatID)
	if err != nil {
		return "", fmt.Errorf("ошибка получения контекста: %w", err)
	}

	// 3. Формируем запрос к Python AI Service
	pyReq := PythonChatRequest{
		ChatID:     chatID,
		NewMessage: userText,
		Context:    contextMsgs,
	}
	reqBytes, _ := json.Marshal(pyReq)

	// 4. Делаем HTTP POST запрос в Python
	resp, err := http.Post(s.aiURL, "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		return "", fmt.Errorf("ошибка запроса к AI сервису: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("AI сервис вернул ошибку: статус %d", resp.StatusCode)
	}

	// 5. Читаем ответ от Gemini
	body, _ := io.ReadAll(resp.Body)
	var pyResp PythonChatResponse
	if err := json.Unmarshal(body, &pyResp); err != nil {
		return "", fmt.Errorf("ошибка парсинга ответа AI: %w", err)
	}

	// 6. Сохраняем ответ ИИ в Postgres
	aiMsg := &domain.Message{
		ChatID:  chatID,
		Role:    "model",
		Content: pyResp.Response,
	}
	if err := s.pgRepo.SaveMessage(ctx, aiMsg); err != nil {
		return "", fmt.Errorf("ошибка сохранения ответа ИИ: %w", err)
	}

	// 7. Обновляем контекст в Redis (добавляем туда новые сообщения)
	contextMsgs = append(contextMsgs, *userMsg, *aiMsg)
	if len(contextMsgs) > 10 { // Храним только последние 10 сообщений
		contextMsgs = contextMsgs[len(contextMsgs)-10:]
	}
	
	if err := s.rdRepo.SaveContext(ctx, chatID, contextMsgs); err != nil {
		// Ошибка кэша не критична для ответа юзеру, просто залогируем 
		fmt.Printf("Внимание: не удалось обновить кэш Redis: %v\n", err)
	}

	return pyResp.Response, nil
}

func (s *ChatService) CreateNewChat(ctx context.Context) (string, error) {
    return s.pgRepo.CreateChat(ctx)
}

