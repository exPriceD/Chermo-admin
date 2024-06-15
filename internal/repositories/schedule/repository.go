package schedule

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/exPriceD/Chermo-admin/internal/entities"
	"github.com/exPriceD/Chermo-admin/internal/models"
	"github.com/jmoiron/sqlx"
	"time"
)

type Repository struct {
	db *sqlx.DB
}

func NewScheduleRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetEventDates(eventID int) ([]string, error) {
	query := "SELECT event_date FROM event_schedule WHERE event_id = $1"
	rows, err := r.db.Query(query, eventID)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	var dates []string
	for rows.Next() {
		var eventDate time.Time
		if err := rows.Scan(&eventDate); err != nil {
			return nil, err
		}
		formattedDate := eventDate.Format("02-01-2006")
		dates = append(dates, formattedDate)
	}
	return dates, nil
}

func (r *Repository) GetTimeSlots(eventID int, eventDate string) ([]map[string]interface{}, error) {
	parsedDate, err := time.Parse("02-01-2006", eventDate)
	if err != nil {
		return nil, err
	}

	query := `SELECT et.id, et.start_time, et.end_time, et.total_slots, et.available_slots
                 FROM event_timeslots et
                 JOIN event_schedule es ON et.schedule_id = es.id
                 WHERE es.event_id = $1 AND es.event_date = $2`
	rows, err := r.db.Query(query, eventID, parsedDate.Format("2006-01-02"))
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	var timeSlots []map[string]interface{}
	for rows.Next() {
		var id, totalSlots, availableSlots int
		var startTime, endTime time.Time
		if err := rows.Scan(&id, &startTime, &endTime, &totalSlots, &availableSlots); err != nil {
			return nil, err
		}
		timeSlot := map[string]interface{}{
			"id":              id,
			"start_time":      startTime.Format("15:04"),
			"end_time":        endTime.Format("15:04"),
			"total_slots":     totalSlots,
			"available_slots": availableSlots,
		}
		timeSlots = append(timeSlots, timeSlot)
	}
	return timeSlots, nil
}

func (r *Repository) GetScheduleByEventID(eventID int) ([]entities.Schedule, error) {
	rows, err := r.db.Queryx(`
        SELECT es.event_date, et.id, et.start_time, et.end_time, et.total_slots, et.available_slots
        FROM event_schedule es
        JOIN event_timeslots et ON es.id = et.schedule_id
        WHERE es.event_id = $1
        ORDER BY es.event_date, et.start_time`, eventID)
	if err != nil {
		return nil, err
	}
	defer func(rows *sqlx.Rows) {
		_ = rows.Close()
	}(rows)

	scheduleMap := make(map[string][]entities.TimeSlot)
	for rows.Next() {
		var eventDate time.Time
		var startTime, endTime time.Time
		var timeSlot entities.TimeSlot
		var timeSlotID int
		err := rows.Scan(&eventDate, &timeSlotID, &startTime, &endTime, &timeSlot.TotalSlots, &timeSlot.AvailableSlots)
		if err != nil {
			return nil, err
		}
		dateStr := eventDate.Format("2006-01-02")
		timeSlot.ID = timeSlotID
		timeSlot.StartTime = startTime.Format("15:04")
		timeSlot.EndTime = endTime.Format("15:04")
		scheduleMap[dateStr] = append(scheduleMap[dateStr], timeSlot)
	}

	var schedules []entities.Schedule
	for date, slots := range scheduleMap {
		schedules = append(schedules, entities.Schedule{
			EventDate: date,
			TimeSlots: slots,
		})
	}

	return schedules, nil
}

func (r *Repository) CreateSchedule(startDate time.Time, endDate time.Time, req models.ScheduleRequest) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	for date := startDate; !date.After(endDate); date = date.AddDate(0, 0, 1) {
		dayOfWeek := date.Weekday().String()

		if timeSlots, ok := req.TimeSlots[dayOfWeek]; ok {
			var scheduleID int
			err := tx.QueryRowx(`
                SELECT id FROM event_schedule
                WHERE event_id = $1 AND event_date = $2`, req.EventID, date).Scan(&scheduleID)

			if err != nil && !errors.Is(err, sql.ErrNoRows) {
				_ = tx.Rollback()
				return err
			}

			if errors.Is(err, sql.ErrNoRows) {
				// Если расписание не существует, создаем новое
				err = tx.QueryRowx(`
                    INSERT INTO event_schedule (event_id, event_date)
                    VALUES ($1, $2)
                    RETURNING id`, req.EventID, date).Scan(&scheduleID)
				if err != nil {
					_ = tx.Rollback()
					return err
				}
			} else {
				// Если расписание уже существует, удаляем существующие временные слоты
				_, err = tx.Exec(`
                    DELETE FROM event_timeslots
                    WHERE schedule_id = $1`, scheduleID)
				if err != nil {
					_ = tx.Rollback()
					return err
				}
			}

			for _, slot := range timeSlots {
				startTime, err := time.Parse("15:04", slot.StartTime)
				if err != nil {
					_ = tx.Rollback()
					return err
				}

				endTime := startTime.Add(time.Duration(req.Duration) * time.Minute).Format("15:04")

				_, err = tx.Exec(`
                    INSERT INTO event_timeslots (schedule_id, start_time, end_time, total_slots, available_slots)
                    VALUES ($1, $2, $3, $4, $5)`,
					scheduleID, slot.StartTime, endTime, slot.Slots, slot.Slots)
				if err != nil {
					_ = tx.Rollback()
					return err
				}
			}
		}
	}

	err = tx.Commit()

	return nil
}

func (r *Repository) TimeSlotExists(timeslotID int) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM event_timeslots WHERE id = $1)"
	err := r.db.QueryRow(query, timeslotID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *Repository) IsTimeSlotAvailable(timeslotID int) (bool, error) {
	var availableSlots int
	err := r.db.QueryRow("SELECT available_slots FROM event_timeslots WHERE id = $1", timeslotID).Scan(&availableSlots)
	if err != nil {
		return false, err
	}
	return availableSlots > 0, nil
}

func (r *Repository) RegisterVisitor(timeslotID, visitorID int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	var availableSlots int
	err = tx.QueryRow("SELECT available_slots FROM event_timeslots WHERE id = $1 FOR UPDATE", timeslotID).Scan(&availableSlots)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	if availableSlots <= 0 {
		_ = tx.Rollback()
		return fmt.Errorf("no available slots")
	}

	// Регистрация посетителя с обработкой конфликта
	_, err = tx.Exec(`
        INSERT INTO event_registrations (timeslot_id, visitor_id) 
        VALUES ($1, $2)
        ON CONFLICT (timeslot_id, visitor_id) DO NOTHING
    `, timeslotID, visitorID)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	// Обновление количества доступных слотов, только если регистрация была успешной
	_, err = tx.Exec("UPDATE event_timeslots SET available_slots = available_slots - 1 WHERE id = $1 AND available_slots > 0", timeslotID)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *Repository) UpdateRegistrationStatus(visitorID, timeslotID int, isConfirmed bool) error {
	query := "UPDATE event_registrations SET is_confirmed = $1 WHERE visitor_id = $2 AND timeslot_id = $3"
	result, err := r.db.Exec(query, isConfirmed, visitorID, timeslotID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no registration found with")
	}

	return nil
}
