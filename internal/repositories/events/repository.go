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
	err := r.db.Select(&events, "SELECT id, title, description, image_url, museum_id FROM events")
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (r *Repository) GetEventByID(eventID int) (entities.Event, error) {
	var event entities.Event
	err := r.db.Get(&event, "SELECT id, title, description, image_url, museum_id FROM events WHERE id = $1", eventID)
	if err != nil {
		return event, err
	}

	return event, nil
}

func (r *Repository) GetEventsByMuseum(museumID int) ([]entities.Event, error) {
	var events []entities.Event
	err := r.db.Select(&events, "SELECT id, title, description, image_url, museum_id FROM events WHERE museum_id = $1", museumID)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (r *Repository) InsertEvent(event entities.Event) error {
	query := `
        INSERT INTO events (title, description, image_url, museum_id)
        VALUES ($1, $2, $3, $4)
    `

	_, err := r.db.Exec(query, event.Title, event.Description, event.ImageURL, event.MuseumID)
	return err
}
