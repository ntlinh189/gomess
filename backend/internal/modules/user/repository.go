package user

import (
	"gomess/internal/database"
	"gomess/internal/models"
)

type RepositoryInterface interface {
	GetByID(id int64) (*models.User, error)
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
		SELECT
			id,
			provider,
			provider_id,
			email,
			name,
			avatar
		FROM users
		WHERE id = ?
		LIMIT 1
	`

	err := r.db.GetDB().QueryRow(query, id).Scan(
		&user.ID,
		&user.Provider,
		&user.ProviderID,
		&user.Email,
		&user.Name,
		&user.Avatar,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
