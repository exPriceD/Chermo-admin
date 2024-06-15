package visitors

import (
	"database/sql"
	"fmt"
	"github.com/exPriceD/Chermo-admin/internal/entities"
	"github.com/jmoiron/sqlx"
	"time"
)

type Repository struct {
	db *sqlx.DB
}

func NewVisitorsRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) AddVisitor(visitor entities.Visitor) (int, error) {
	var visitorID int
	err := r.db.QueryRow(`
        INSERT INTO visitors (first_name, last_name, patronymic, phone, email)
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT (email) DO UPDATE SET
            first_name = EXCLUDED.first_name,
            last_name = EXCLUDED.last_name,
            patronymic = EXCLUDED.patronymic,
            phone = EXCLUDED.phone
        RETURNING id`, visitor.FirstName, visitor.LastName, visitor.Patronymic, visitor.Phone, visitor.Email).Scan(&visitorID)
	if err != nil {
		return 0, err
	}
	return visitorID, nil
}

func (r *Repository) VisitorExists(visitorID int) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM visitors WHERE id = $1)"
	err := r.db.QueryRow(query, visitorID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *Repository) GetVisitors(eventID int, eventDate, startTime string) ([]map[string]interface{}, error) {
	parsedDate, err := time.Parse("02-01-2006", eventDate)
	if err != nil {
		return nil, err
	}
	fmt.Println(parsedDate.Format("2006-01-02"))

	parsedTime, err := time.Parse("15:04", startTime)
	if err != nil {
		return nil, err
	}
	fmt.Println(parsedTime)
	query := `SELECT v.id, v.first_name, v.last_name, v.patronymic, v.phone, v.email
                 FROM visitors v
                 JOIN event_registrations er ON v.id = er.visitor_id
                 JOIN event_timeslots et ON er.timeslot_id = et.id
                 JOIN event_schedule es ON et.schedule_id = es.id
                 WHERE es.event_id = $1 AND es.event_date = $2 AND et.start_time = $3`
	rows, err := r.db.Query(query, eventID, parsedDate.Format("2006-01-02"), parsedTime)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	var visitors []map[string]interface{}
	for rows.Next() {
		var id int
		var firstName, lastName, patronymic, phone, email string
		if err := rows.Scan(&id, &firstName, &lastName, &patronymic, &phone, &email); err != nil {
			return nil, err
		}
		visitor := map[string]interface{}{
			"id":         id,
			"first_name": firstName,
			"last_name":  lastName,
			"patronymic": patronymic,
			"phone":      phone,
			"email":      email,
		}
		visitors = append(visitors, visitor)
	}
	return visitors, nil
}
