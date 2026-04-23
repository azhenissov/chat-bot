package api

import (
	"net/http"

	"chat-core/internal/service"

	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	chatService *service.ChatService
}

func NewChatHandler(s *service.ChatService) *ChatHandler {
	return &ChatHandler{chatService: s}
}


func (h *ChatHandler) StartChat(c *gin.Context) {
	// Request.Context() передает контекст HTTP-запроса дальше в базу
	chatID, err := h.chatService.CreateNewChat(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot create chat", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"chat_id": chatID})
}

type MessageRequest struct {
	ChatID string `json:"chat_id" binding:"required"`
	Message string `json:"message" binding:"required"`
}

func (h *ChatHandler) SendMessage(c *gin.Context) {
	var req MessageRequest

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	reply, err := h.chatService.ProcessMessage(c.Request.Context(), req.ChatID, req.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot process message", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"chat_id": req.ChatID,
		"reply": reply,
	})

	
	
}