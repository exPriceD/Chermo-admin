package entities

type Event struct {
	ID          int    `json:"id" db:"id"`
	Title       string `json:"title" db:"title"`
	Description string `json:"description" db:"description"`
	ImageURL    string `json:"image_url" db:"image_url"`
}
