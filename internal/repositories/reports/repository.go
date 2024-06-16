package reports

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewReportsRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetVisitorReport(eventID string) (*sql.Rows, error) {
	query := `
        SELECT
            er.registration_date, v.first_name, v.phone
        FROM event_registrations er
        JOIN visitors v ON er.visitor_id = v.id
        JOIN event_timeslots et ON er.timeslot_id = et.id
        JOIN event_schedule es ON et.schedule_id = es.id
        JOIN events e ON es.event_id = e.id
        WHERE e.id = $1 AND er.is_confirmed = true
    `
	return r.db.Query(query, eventID)
}

func (r *Repository) GetEventReport() (*sql.Rows, error) {
	query := `
        SELECT
            m.name, es.event_date, et.start_time, et.end_time, e.title, et.total_slots - COUNT(er.id) AS available_slots
        FROM events e
        JOIN event_schedule es ON e.id = es.event_id
        JOIN event_timeslots et ON es.id = et.schedule_id
        JOIN museums m ON e.museum_id = m.id
        LEFT JOIN event_registrations er ON et.id = er.timeslot_id AND er.is_confirmed = true
        GROUP BY m.name, es.event_date, et.start_time, et.end_time, e.title, et.total_slots
    `
	return r.db.Query(query)
}
