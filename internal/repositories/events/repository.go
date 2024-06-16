package events

import (
	"database/sql"
	"errors"
	"github.com/exPriceD/Chermo-admin/internal/entities"
	"github.com/jmoiron/sqlx"
	"time"
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

type Event struct {
	Title     string
	Museum    string
	Date      time.Time
	StartTime time.Time
	EndTime   time.Time
}

func (r *Repository) GetEventByTimeslotID(timeslotID int) (*Event, error) {
	query := `
        SELECT e.title, m.name, es.event_date, et.start_time, et.end_time
        FROM events e
        JOIN event_schedule es ON e.id = es.event_id
        JOIN event_timeslots et ON es.id = et.schedule_id
        JOIN museums m ON e.museum_id = m.id
        WHERE et.id = $1
    `

	var event Event
	err := r.db.QueryRow(query, timeslotID).Scan(&event.Title, &event.Museum, &event.Date, &event.StartTime, &event.EndTime)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &event, nil
}
