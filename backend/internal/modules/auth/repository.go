package auth

import (
	"database/sql"
	"gomess/internal/database"
	"gomess/internal/models"
)

type RepositoryInterface interface {
	FindByProvider(provider, providerID string) (*models.User, error)
	Create(user *models.User) error
}

type Repository struct {
	db database.DatabaseInterface
}

func NewRepository(db database.DatabaseInterface) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindByProvider(provider, providerID string) (*models.User, error) {
	query := `
	SELECT id, provider, provider_id, email, name, avatar
	FROM users
	WHERE provider = ? AND provider_id = ?
	LIMIT 1
	`

	row := r.db.GetDB().QueryRow(query, provider, providerID)

	var user models.User

	err := row.Scan(
		&user.ID,
		&user.Provider,
		&user.ProviderID,
		&user.Email,
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

func (r *Repository) Create(user *models.User) error {
	query := `
	INSERT INTO users(provider, provider_id, email, name, avatar)
	VALUES (?, ?, ?, ?, ?)
	`

	result, err := r.db.GetDB().Exec(
		query,
		user.Provider,
		user.ProviderID,
		user.Email,
		user.Name,
		user.Avatar,
	)

	if err != nil {
		return err
	}

	id, _ := result.LastInsertId()

	user.ID = id

	return nil
}
