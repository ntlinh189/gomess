package user

import "gomess/internal/database"

type RepositoryInterface interface {
}

type Repository struct {
	db database.DatabaseInterface
}

func NewRepository(db database.DatabaseInterface) *Repository {
	return &Repository{db: db}
}
