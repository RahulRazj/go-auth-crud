package repository

import "github.com/RahulRazj/go-crud-auth/internal/database"

type BaseRepository struct {
	db *database.Database
}

func NewBaseRepository(db *database.Database) BaseRepository {
	return BaseRepository{db: db}
}
