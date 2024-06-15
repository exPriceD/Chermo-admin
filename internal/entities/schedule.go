package entities

type TimeSlot struct {
	StartTime      string `json:"start_time"`
	EndTime        string `json:"end_time"`
	TotalSlots     int    `json:"total_slots"`
	AvailableSlots int    `json:"available_slots"`
}

type Schedule struct {
	EventDate string     `json:"event_date"`
	TimeSlots []TimeSlot `json:"time_slots"`
}
