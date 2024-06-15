package auth

import (
	"github.com/exPriceD/Chermo-admin/internal/entities"
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

func (r *Repository) GetUser(username string) (entities.ReceivedUser, error) {
	var receivedUser entities.ReceivedUser

	err := r.db.Get(&receivedUser, "SELECT id, role, password, museum_id FROM users WHERE username = $1", username)
	if err != nil {
		return receivedUser, err
	}

	return receivedUser, nil
}

func (r *Repository) InsertUser(u *models.User) error {
	query := `
		INSERT INTO users (username, password, role, museum_id)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.Exec(query, u.Username, u.Password, u.Role, u.MuseumID)
	return err
}
