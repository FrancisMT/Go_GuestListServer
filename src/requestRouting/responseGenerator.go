package requestRouting

import (
	"guestListChallenge/src/database"
)

// CreateAddGuestResponse Creates a response for "add a guest to the guest list" requests
//
// A struct with the appropriate fields and json tags is used
func CreateAddGuestResponse(guest database.GuestList) interface{} {
	return struct {
		Name string `json:"name"`
	}{Name: guest.Name}
}

// CreateGetGuestListResponse Creates a response for "get the guest list" requests
//
// A struct with the appropriate fields and json tags is used
func CreateGetGuestListResponse(guestList []database.GuestList) interface{} {

	// Guest data to send in the response
	type guestData struct {
		Name               string `json:"name"`
		Table              int    `json:"table"`
		AccompanyingGuests int    `json:"accompanying_guests"`
	}

	// Populate guest data array
	guestDataArray := make([]guestData, 0, len(guestList))
	for _, guest := range guestList {
		guestDataArray = append(guestDataArray, guestData{guest.Name, guest.Table, guest.AccompanyingGuests})
	}

	return struct {
		Guests []guestData `json:"guests"`
	}{Guests: guestDataArray}
}

// CreateCheckInGuestResponse Creates a response for "guest arrives to the party" requests
//
// A struct with the appropriate fields and json tags is used
func CreateCheckInGuestResponse(guest database.GuestList) interface{} {
	return struct {
		Name string `json:"name"`
	}{Name: guest.Name}
}

// CreateGetArrivedGuestsResponse Creates a response for "get list of guests that have arrived to the party" requests
//
// A struct with the appropriate fields and json tags is used
func CreateGetArrivedGuestsResponse(guestList []database.GuestList) interface{} {

	// Guest data to send in the response
	type guestData struct {
		Name               string `json:"name"`
		AccompanyingGuests int    `json:"accompanying_guests"`
		TimeArrived        string `json:"time_arrived"`
	}

	// Populate guest data array
	guestDataArray := make([]guestData, 0, len(guestList))
	for _, guest := range guestList {
		guestDataArray = append(guestDataArray, guestData{guest.Name, guest.AccompanyingGuests, guest.TimeArrived})
	}

	return struct {
		Guests []guestData `json:"guests"`
	}{Guests: guestDataArray}
}

// CreateGetNumberOfEmptySeatsResponse Creates a response for "get the number of empty seats" requests
//
// A struct with the appropriate fields and json tags is used
func CreateGetNumberOfEmptySeatsResponse(seatsEmpty int) interface{} {
	return struct {
		SeatsEmpty int `json:"seats_empty"`
	}{SeatsEmpty: seatsEmpty}
}
