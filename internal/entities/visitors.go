package entities

type Visitor struct {
	ID         int    `json:"id" db:"id"`
	FirstName  string `json:"first_name" db:"first_name"`
	LastName   string `json:"last_name" db:"last_name"`
	Patronymic string `json:"patronymic" db:"patronymic"`
	Phone      string `json:"phone" db:"phone"`
	Email      string `json:"email" db:"email"`
}
