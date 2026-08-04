package message

import (
	"database/sql"
	"gomess/internal/database"
	"gomess/internal/models"
	"strings"
)

type RepositoryInterface interface {
	GetByID(id int64) (*models.Message, error)
	GetHistory(userID, friendID int64, beforeID int64, limit int) ([]models.Message, error)
	GetAttachmentsByMessageIDs(messageIDs []int64) (map[int64][]models.MessageAttachment, error)
	CreateDeletion(messageID, userID int64) error
	RevokeMessage(id int64) error
	GetAllObjectKeys() ([]string, error)

	// Transaction
	CreateTx(tx *sql.Tx, message *models.Message) error
	CreateAttachmentsTx(tx *sql.Tx, messageID int64, attachments []models.MessageAttachment) error
	DeleteConversationTx(tx *sql.Tx, user1ID, user2ID int64) error
}

type Repository struct {
	db database.DatabaseInterface
}

func NewRepository(db database.DatabaseInterface) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetByID(id int64) (*models.Message, error) {
	query := `
	SELECT id, sender_id, receiver_id, content, created_at, revoked_at
	FROM messages
	WHERE id = ?
	LIMIT 1
	`

	var m models.Message
	err := r.db.GetDB().QueryRow(query, id).Scan(
		&m.ID, &m.SenderID, &m.ReceiverID, &m.Content, &m.CreatedAt, &m.RevokedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &m, nil
}

func (r *Repository) GetHistory(userID, friendID int64, beforeID int64, limit int) ([]models.Message, error) {
	query := `
	SELECT m.id, m.sender_id, m.receiver_id, m.content, m.created_at, m.revoked_at
	FROM messages m
	WHERE (
		(m.sender_id = ? AND m.receiver_id = ?)
		OR (m.sender_id = ? AND m.receiver_id = ?)
	)
	AND NOT EXISTS (
		SELECT 1 FROM message_deletions d
		WHERE d.message_id = m.id AND d.user_id = ?
	)
	`

	args := []any{userID, friendID, friendID, userID, userID}

	if beforeID > 0 {
		query += " AND m.id < ?"
		args = append(args, beforeID)
	}

	query += " ORDER BY m.id DESC LIMIT ?"
	args = append(args, limit)

	rows, err := r.db.GetDB().Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := make([]models.Message, 0, limit)

	for rows.Next() {
		var m models.Message
		if err := rows.Scan(&m.ID, &m.SenderID, &m.ReceiverID, &m.Content, &m.CreatedAt, &m.RevokedAt); err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}

	return messages, nil
}

func (r *Repository) GetAttachmentsByMessageIDs(messageIDs []int64) (map[int64][]models.MessageAttachment, error) {
	result := make(map[int64][]models.MessageAttachment)

	if len(messageIDs) == 0 {
		return result, nil
	}

	placeholders := make([]string, len(messageIDs))
	args := make([]any, len(messageIDs))
	for i, id := range messageIDs {
		placeholders[i] = "?"
		args[i] = id
	}

	query := `
	SELECT id, message_id, type, object_key, file_name, mime_type, size_bytes, created_at
	FROM message_attachments
	WHERE message_id IN (` + strings.Join(placeholders, ",") + `)
	`

	rows, err := r.db.GetDB().Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var a models.MessageAttachment
		err := rows.Scan(
			&a.ID, &a.MessageID, &a.Type, &a.ObjectKey,
			&a.FileName, &a.MimeType, &a.SizeBytes, &a.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		result[a.MessageID] = append(result[a.MessageID], a)
	}

	return result, nil
}

func (r *Repository) CreateDeletion(messageID, userID int64) error {
	query := `
	INSERT IGNORE INTO message_deletions (message_id, user_id)
	VALUES (?, ?)
	`
	_, err := r.db.GetDB().Exec(query, messageID, userID)
	return err
}

func (r *Repository) RevokeMessage(id int64) error {
	query := `
	UPDATE messages
	SET revoked_at = CURRENT_TIMESTAMP
	WHERE id = ?
	`
	_, err := r.db.GetDB().Exec(query, id)
	return err
}

func (r *Repository) GetAllObjectKeys() ([]string, error) {
	query := `SELECT object_key FROM message_attachments`

	rows, err := r.db.GetDB().Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var keys []string
	for rows.Next() {
		var key string
		if err := rows.Scan(&key); err != nil {
			return nil, err
		}
		keys = append(keys, key)
	}

	return keys, nil
}

func (r *Repository) CreateTx(tx *sql.Tx, message *models.Message) error {
	query := `
	INSERT INTO messages (sender_id, receiver_id, content)
	VALUES (?, ?, ?)
	`

	result, err := tx.Exec(query, message.SenderID, message.ReceiverID, message.Content)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	message.ID = id

	return nil
}

func (r *Repository) CreateAttachmentsTx(tx *sql.Tx, messageID int64, attachments []models.MessageAttachment) error {
	if len(attachments) == 0 {
		return nil
	}

	query := `
	INSERT INTO message_attachments (message_id, type, object_key, file_name, mime_type, size_bytes)
	VALUES (?, ?, ?, ?, ?, ?)
	`

	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, a := range attachments {
		if _, err := stmt.Exec(messageID, a.Type, a.ObjectKey, a.FileName, a.MimeType, a.SizeBytes); err != nil {
			return err
		}
	}

	return nil
}

func (r *Repository) DeleteConversationTx(tx *sql.Tx, user1ID, user2ID int64) error {
	query := `
	DELETE FROM messages
	WHERE (sender_id = ? AND receiver_id = ?)
	   OR (sender_id = ? AND receiver_id = ?)
	`
	_, err := tx.Exec(query, user1ID, user2ID, user2ID, user1ID)
	return err
}
