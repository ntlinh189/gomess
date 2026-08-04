package friend

import (
	"database/sql"
	"gomess/internal/database"
	"gomess/internal/models"
)

type RepositoryInterface interface {
	IsFriend(user1ID, user2ID int64) (bool, error)
	FindPendingRequest(senderID, receiverID int64) (*models.FriendRequest, error)
	CreateRequest(request *models.FriendRequest) error
	GetReceivedRequests(userID int64) ([]models.FriendRequest, error)
	GetSentRequests(userID int64) ([]models.FriendRequest, error)
	GetRequestByID(id int64) (*models.FriendRequest, error)
	UpdateRequestStatus(id int64, status string) error
	CreateFriend(user1ID, user2ID int64) error
	DeleteFriend(user1ID, user2ID int64) error
	GetFriends(userID int64) ([]models.User, error)

	// Transaction
	UpdateRequestStatusTx(tx *sql.Tx, id int64, status string) error
	CreateFriendTx(tx *sql.Tx, user1ID, user2ID int64) error
	DeleteFriendTx(tx *sql.Tx, user1ID, user2ID int64) error
}

type Repository struct {
	db database.DatabaseInterface
}

func NewRepository(db database.DatabaseInterface) *Repository {
	return &Repository{db: db}
}

func (r *Repository) IsFriend(user1ID, user2ID int64) (bool, error) {
	if user1ID > user2ID {
		user1ID, user2ID = user2ID, user1ID
	}

	query := `
	SELECT EXISTS(
		SELECT 1
		FROM friends
		WHERE user1_id = ? AND user2_id = ?
	)
	`

	var exists bool
	err := r.db.GetDB().QueryRow(query, user1ID, user2ID).Scan(&exists)

	return exists, err
}

func (r *Repository) FindPendingRequest(senderID, receiverID int64) (*models.FriendRequest, error) {
	query := `
	SELECT
		id, sender_id, receiver_id, status, created_at, updated_at
	FROM friend_requests
	WHERE (
		(sender_id = ? AND receiver_id = ?)
		OR (sender_id = ? AND receiver_id = ?)
	) AND status = 'pending'
	LIMIT 1
	`

	row := r.db.GetDB().QueryRow(query, senderID, receiverID, receiverID, senderID)
	var req models.FriendRequest
	err := row.Scan(
		&req.ID,
		&req.SenderID,
		&req.ReceiverID,
		&req.Status,
		&req.CreatedAt,
		&req.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &req, nil
}

func (r *Repository) CreateRequest(request *models.FriendRequest) error {
	query := `
	INSERT INTO friend_requests (sender_id, receiver_id, status)
	VALUES (?, ?, ?)
	`
	result, err := r.db.GetDB().Exec(query, request.SenderID, request.ReceiverID, request.Status)

	if err != nil {
		return err
	}

	id, _ := result.LastInsertId()
	request.ID = id
	return nil
}

func (r *Repository) GetReceivedRequests(userID int64) ([]models.FriendRequest, error) {
	query := `
	SELECT id, sender_id, receiver_id, status, created_at, updated_at
	FROM friend_requests
	WHERE receiver_id = ?
	`
	rows, err := r.db.GetDB().Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	requests := make([]models.FriendRequest, 0)

	for rows.Next() {
		var req models.FriendRequest
		err := rows.Scan(
			&req.ID,
			&req.SenderID,
			&req.ReceiverID,
			&req.Status,
			&req.CreatedAt,
			&req.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		requests = append(requests, req)
	}
	return requests, nil
}

func (r *Repository) GetSentRequests(userID int64) ([]models.FriendRequest, error) {
	query := `
	SELECT id, sender_id, receiver_id, status, created_at, updated_at
	FROM friend_requests
	WHERE sender_id = ?
	`
	rows, err := r.db.GetDB().Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	requests := make([]models.FriendRequest, 0)

	for rows.Next() {
		var req models.FriendRequest
		err := rows.Scan(
			&req.ID,
			&req.SenderID,
			&req.ReceiverID,
			&req.Status,
			&req.CreatedAt,
			&req.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		requests = append(requests, req)
	}
	return requests, nil
}

func (r *Repository) GetRequestByID(id int64) (*models.FriendRequest, error) {
	query := `
	SELECT id, sender_id, receiver_id, status, created_at, updated_at
	FROM friend_requests
	WHERE id = ?
	LIMIT 1
	`
	row := r.db.GetDB().QueryRow(query, id)
	var req models.FriendRequest
	err := row.Scan(
		&req.ID,
		&req.SenderID,
		&req.ReceiverID,
		&req.Status,
		&req.CreatedAt,
		&req.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &req, nil
}

func (r *Repository) UpdateRequestStatus(id int64, status string) error {
	query := `
	UPDATE friend_requests
	SET status = ?
	WHERE id = ?
	`
	_, err := r.db.GetDB().Exec(query, status, id)
	return err
}

func (r *Repository) CreateFriend(user1ID, user2ID int64) error {
	if user1ID > user2ID {
		user1ID, user2ID = user2ID, user1ID
	}

	query := `
	INSERT INTO friends(user1_id, user2_id)
	VALUES(?, ?)
	`

	_, err := r.db.GetDB().Exec(query, user1ID, user2ID)
	return err
}

func (r *Repository) DeleteFriend(user1ID, user2ID int64) error {
	if user1ID > user2ID {
		user1ID, user2ID = user2ID, user1ID
	}

	query := `
	DELETE FROM friends
	WHERE user1_id = ? AND user2_id = ?
	`

	_, err := r.db.GetDB().Exec(query, user1ID, user2ID)
	return err
}

func (r *Repository) GetFriends(userID int64) ([]models.User, error) {
	query := `
	SELECT u.id, u.provider, u.provider_id, u.account, u.name, u.avatar
	FROM friends f
	JOIN users u 
	ON u.id = 
		CASE
			WHEN f.user1_id = ? THEN f.user2_id ELSE f.user1_id
		END
	WHERE f.user1_id = ? OR f.user2_id = ?
	ORDER BY u.name
	`

	rows, err := r.db.GetDB().Query(query, userID, userID, userID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := make([]models.User, 0)

	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID,
			&user.Provider,
			&user.ProviderID,
			&user.Account,
			&user.Name,
			&user.Avatar,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *Repository) UpdateRequestStatusTx(tx *sql.Tx, id int64, status string) error {
	query := `
	UPDATE friend_requests
	SET status = ?
	WHERE id = ?
	`

	_, err := tx.Exec(query, status, id)
	return err
}

func (r *Repository) CreateFriendTx(tx *sql.Tx, user1ID, user2ID int64) error {
	if user1ID > user2ID {
		user1ID, user2ID = user2ID, user1ID
	}
	
	query := `
	INSERT INTO friends(user1_id, user2_id)
	VALUES(?, ?)
	`
	_, err := tx.Exec(query, user1ID, user2ID)
	return err
}

func (r *Repository) DeleteFriendTx(tx *sql.Tx, user1ID, user2ID int64) error {
	if user1ID > user2ID {
		user1ID, user2ID = user2ID, user1ID
	}

	query := `
	DELETE FROM friends
	WHERE user1_id = ? AND user2_id = ?
	`

	_, err := tx.Exec(query, user1ID, user2ID)
	return err
}