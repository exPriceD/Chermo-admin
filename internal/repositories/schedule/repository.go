package schedule

import (
	"database/sql"
	"errors"
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

func (r *Repository) GetScheduleByEventID(eventID int) ([]entities.Schedule, error) {
	rows, err := r.db.Queryx(`
        SELECT es.event_date, et.start_time, et.end_time, et.total_slots, et.available_slots
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
		err := rows.Scan(&eventDate, &startTime, &endTime, &timeSlot.TotalSlots, &timeSlot.AvailableSlots)
		if err != nil {
			return nil, err
		}
		dateStr := eventDate.Format("2006-01-02")
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
