package user

import (
	"database/sql"
	"gomess/internal/database"
	"gomess/internal/models"
)

type RepositoryInterface interface {
	GetByID(id int64) (*models.User, error)
	Search(userID int64, provider string, keyword string, skip int, limit int) ([]*models.User, error)
	ExistsUser(userID int64) (bool, error)
	DeleteUser(userID int64) error
}

type Repository struct {
	db database.DatabaseInterface
}

func NewRepository(db database.DatabaseInterface) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetByID(id int64) (*models.User, error) {
	var user models.User

	query := `
	SELECT id, provider, provider_id, account, name, avatar
	FROM users
	WHERE id = ?
	LIMIT 1
	`

	err := r.db.GetDB().QueryRow(query, id).Scan(
		&user.ID,
		&user.Provider,
		&user.ProviderID,
		&user.Account,
		&user.Name,
		&user.Avatar,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) Search(userID int64, provider string, keyword string, skip int, limit int) ([]*models.User, error) {
	query := `
	SELECT u.id, u.provider, u.provider_id, u.account, u.name, u.avatar,
		COALESCE(fr.status, '')
	FROM users u
	LEFT JOIN friend_requests fr
		ON fr.sender_id = ? AND fr.receiver_id = u.id AND fr.status = 'pending'
	WHERE u.provider = ? AND LOWER(u.account) LIKE LOWER(?)
	ORDER BY u.name
	LIMIT ? OFFSET ?
	`

	rows, err := r.db.GetDB().Query(query, userID, provider, keyword, limit, skip)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := []*models.User{}

	for rows.Next() {
		var user models.User

		err := rows.Scan(
			&user.ID,
			&user.Provider,
			&user.ProviderID,
			&user.Account,
			&user.Name,
			&user.Avatar,
			&user.RequestStatus,
		)

		if err != nil {
			return nil, err
		}

		users = append(users, &user)
	}

	return users, nil
}

func (r *Repository) ExistsUser(userID int64) (bool, error) {
	query := `
	SELECT EXISTS(
		SELECT 1
		FROM users
		WHERE id = ?
	)
	`
	var exists bool
	err := r.db.GetDB().QueryRow(query, userID).Scan(&exists)

	return exists, err
}

func (r *Repository) DeleteUser(userID int64) error {
	query := `
	DELETE FROM users
	WHERE id = ?
	`

	_, err := r.db.GetDB().Exec(query, userID)
	return err
}
