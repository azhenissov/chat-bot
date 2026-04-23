package repository

import (

	"context"
	"chat-core/internal/domain"
	"github.com/jackc/pgx/v5"
)

type PostgresRepo struct {
	db *pgx.Conn
}

func NewPostgresRepo(db *pgx.Conn) *PostgresRepo {
	return &PostgresRepo{db: db}
}
// gives back UUID of created chat
func (r *PostgresRepo) CreateChat(ctx context.Context) (string, error) {
	var id string
	err := r.db.QueryRow(ctx, `INSERT INTO chats DEFAULT VALUES RETURNING id;`).Scan(&id)
	return id, err
}


// save chat history in database
func (r *PostgresRepo) SaveMessage(ctx context.Context, msg *domain.Message) error {
	query := `
		INSERT INTO messages (chat_id, role, content) 
		VALUES ($1, $2, $3) 
		RETURNING id, created_at
	` 
	err := r.db.QueryRow(ctx, query, msg.ChatID, msg.Role, msg.Content).Scan(&msg.ID, &msg.CreatedAt)
	
	return err
}


// get chat history
func (r *PostgresRepo) GetChatMessages(ctx context.Context, chatID string) ([]domain.Message, error) {
	query := `SELECT id, chat_id, role, content, created_at FROM messages WHERE chat_id = $1 ORDER BY created_at ASC`
	
	rows, err := r.db.Query(ctx, query, chatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []domain.Message
	for rows.Next() {
		var m domain.Message
		if err := rows.Scan(&m.ID, &m.ChatID, &m.Role, &m.Content, &m.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}
	
	return messages, nil
}

