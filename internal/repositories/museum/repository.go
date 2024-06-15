package museum

import (
	"fmt"
	"github.com/exPriceD/Chermo-admin/internal/entities"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewMuseumRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetMuseums() ([]entities.Museum, error) {
	var museums []entities.Museum
	err := r.db.Select(&museums, "SELECT id, name FROM museums")
	if err != nil {
		return nil, err
	}

	return museums, nil
}

func (r *Repository) GetMuseumByName(name string) (entities.Museum, error) {
	var museum entities.Museum
	fmt.Println(name)
	err := r.db.Get(&museum, "SELECT id, name FROM museums WHERE name = $1", name)
	if err != nil {
		return museum, err
	}

	return museum, nil
}
