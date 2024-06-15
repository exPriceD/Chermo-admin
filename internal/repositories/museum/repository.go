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

func (r *Repository) GetMuseumByName(name string) (entities.Museum, error) {
	var museum entities.Museum
	fmt.Println(name)
	err := r.db.Get(&museum, "SELECT id, name FROM museums WHERE name = $1", name)
	if err != nil {
		return museum, err
	}

	return museum, nil
}
