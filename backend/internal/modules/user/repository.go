package user

import (
	"gomess/internal/database"
	"gomess/internal/models"
)

type RepositoryInterface interface {
	GetByID(id int64) (*models.User, error)
	Search(provider string, keyword string, skip int, limit int) ([]*models.User, error)
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

func (r *Repository) Search(provider string, keyword string, skip int, limit int) ([]*models.User, error) {
	query := `
	SELECT
		id,
		provider,
		provider_id,
		email,
		name,
		avatar
	FROM users
	WHERE provider = ? AND LOWER(email) LIKE LOWER(?)
	ORDER BY name
	LIMIT ? OFFSET ?
	`

	rows, err := r.db.GetDB().Query(query, provider, keyword, limit, skip)

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
			&user.Email,
			&user.Name,
			&user.Avatar,
		)

		if err != nil {
			return nil, err
		}

		users = append(users, &user)
	}

	return users, nil
}
