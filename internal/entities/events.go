package entities

type Event struct {
	ID          int    `json:"id" db:"id"`
	Title       string `json:"title" db:"title"`
	Description string `json:"description" db:"description"`
	MuseumID    int    `json:"museum_id" db:"museum_id"`
	ImageURL    string `json:"image_url" db:"image_url"`
}
