package repository

import (
	"database/sql"
	"gomess/internal/model"
)

type UserRepository interface {
	FindByProvider(provider, providerID string) (*model.User, error)
    Create(user *model.User) error
}

type userRepo struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) FindByProvider(provider, providerID string) (*model.User, error) {
	query := `SELECT id, provider, provider_id, email, name, avatar 
				FROM users WHERE provider = ? AND provider_id = ? LIMIT 1`

	row := r.db.QueryRow(query, provider, providerID)

	var u model.User
	err := row.Scan(&u.ID, &u.Provider, &u.ProviderID, &u.Email, &u.Name, &u.Avatar)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepo) Create(user *model.User) error {
	query := `INSERT INTO users (provider, provider_id, email, name, avatar)
				VALUES (?, ?, ?, ?, ?)`

	res, err := r.db.Exec(query, 
		user.Provider, 
		user.ProviderID, 
		user.Email, 
		user.Name, 
		user.Avatar,
	)
	if err != nil {
		return err
	}

	id, _ := res.LastInsertId()
	user.ID = id
	return nil
}