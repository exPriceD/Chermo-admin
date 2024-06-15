package entities

type Visitor struct {
	ID         int    `json:"id" db:"id"`
	FirstName  string `json:"first_name" db:"first_name"`
	LastName   string `json:"last_name" db:"last_name"`
	Patronymic string `json:"patronymic" db:"patronymic"`
	Phone      string `json:"phone" db:"phone"`
	Email      string `json:"email" db:"email"`
}

type RegistrationRequest struct {
	FirstName  string `json:"first_name" binding:"required"`
	LastName   string `json:"last_name" binding:"required"`
	Patronymic string `json:"patronymic" binding:"required"`
	Phone      string `json:"phone" binding:"required"`
	Email      string `json:"email" binding:"required,email"`
	TimeslotID int    `json:"timeslot_id" binding:"required"`
}
