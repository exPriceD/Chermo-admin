package auth

import (
	"github.com/exPriceD/Chermo-admin/internal/models"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewAuthRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetUser(username, password string) (int, error) {
	var id int
	err := r.db.Get(&id, "SELECT id FROM users WHERE username = $1 AND password = $2", username, password)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *Repository) InsertUser(u *models.User) error {
	query := `
		INSERT INTO users (username, password, role, museum_id)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.Exec(query, u.Username, u.Password, u.MuseumID)
	return err
}
