package database

// GuestList Structure representation of the guestlist sql table used in the database
type GuestList struct {
	Name               string `json:"name" gorm:"primary_key"`
	Table              int    `json:"table"`
	AccompanyingGuests int    `json:"accompanying_guests"`
	TimeArrived        string `json:"time_arrived"`
}
