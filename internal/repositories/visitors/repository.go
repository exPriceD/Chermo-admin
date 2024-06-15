package visitors

import (
	"github.com/exPriceD/Chermo-admin/internal/entities"
	"github.com/jmoiron/sqlx"
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
