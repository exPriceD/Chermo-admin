package entities

type User struct {
	ID         int    `json:"id"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Role       string `json:"role"`
	MuseumName string `json:"museum_name"`
}

type ReceivedUser struct {
	ID       int    `db:"id"`
	Password string `db:"password"`
	Role     string `db:"role"`
}
