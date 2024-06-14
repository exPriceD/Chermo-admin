package events

import (
	"github.com/exPriceD/Chermo-admin/internal/entities"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewEventsRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetEvents() ([]entities.Event, error) {
	var events []entities.Event
	err := r.db.Select(&events, "SELECT id, title, description, image_url FROM events")
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (r *Repository) InsertEvent(event entities.Event) error {
	query := `
        INSERT INTO events (title, description, image_url)
        VALUES ($1, $2, $3)
    `

	_, err := r.db.Exec(query, event.Title, event.Description, event.ImageURL)
	return err
}
