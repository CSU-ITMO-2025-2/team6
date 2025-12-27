package user

import (
	db "local-lib/database"

	"main-service/internal/repository"
)

type repo struct {
	db db.Client
}

func NewRepository(db db.Client) repository.UserRepository {
	return &repo{db: db}
}
