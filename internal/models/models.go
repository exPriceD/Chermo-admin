package models

import "github.com/dgrijalva/jwt-go"

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
	MuseumID int    `json:"museum_id"`
}

type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	MuseumID int    `json:"museum_id"`
	jwt.StandardClaims
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Event struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	MuseumID    int    `json:"museum_id"`
	ImageURL    string `json:"image_url"`
}

type TimeSlot struct {
	StartTime string `json:"start_time"`
	Slots     int    `json:"slots"`
}

type ScheduleRequest struct {
	EventID   int                   `json:"event_id"`
	StartDate string                `json:"start_date"`
	EndDate   string                `json:"end_date"`
	Duration  int                   `json:"duration"`
	TimeSlots map[string][]TimeSlot `json:"time_slots"`
}
