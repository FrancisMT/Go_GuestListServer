package requestRouting

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"guestListChallenge/src/database"
	"guestListChallenge/src/utils"
	"net/http"
)

// encodeResponse Encodes an http response
func encodeResponse(response http.ResponseWriter, reply interface{}) {
	response.Header().Set("Content-Type", "application/json")
	encoderError := json.NewEncoder(response).Encode(reply)
	if encoderError != nil {
		fmt.Println(encoderError.Error())
	}
}

// decodeRequest Decodes an http request and stores the decoded data in a database.GuestList variable
func decodeRequest(request *http.Request) (guest database.GuestList) {
	decoderError := json.NewDecoder(request.Body).Decode(&guest)
	if decoderError != nil {
		fmt.Println(decoderError.Error())
	}
	return
}

// addGuest Processes the request to add a guest to the guest list
//
// An error is reported if the number of accompanying guests is larger than the table capacity
func addGuest(response http.ResponseWriter, request *http.Request) {
	if database.Connector == nil {
		fmt.Println("Database unreachable")
		return
	}

	if request == nil {
		fmt.Println("Null request")
		return
	}

	var requestReply interface{}

	guest := decodeRequest(request)

	// Check table capacity
	if guest.AccompanyingGuests > guest.Table {
		requestReply = "Guest will no be added to the guest list: guest's table cannot hold so many people."
		fmt.Println(requestReply)
	} else {

		// Setup guest data
		guest.Name = mux.Vars(request)["name"]
		guest.TimeArrived = ""

		// Add guest data to database
		database.Connector.Create(&guest)

		requestReply = CreateAddGuestResponse(guest)
	}

	encodeResponse(response, requestReply)
}

// getGuestList Processes the request to get the guest list
func getGuestList(response http.ResponseWriter, _ *http.Request) {
	if database.Connector == nil {
		fmt.Println("Database unreachable")
		return
	}

	var guestList []database.GuestList
	database.Connector.Find(&guestList)
	encodeResponse(response, CreateGetGuestListResponse(guestList))
}

// checkInGuest Processes the request that happens when a guest arrives to the party
//
// An error is reported if the number of accompanying guests is larger than the table capacity
func checkInGuest(response http.ResponseWriter, request *http.Request) {
	if database.Connector == nil {
		fmt.Println("Database unreachable")
		return
	}

	if request == nil {
		fmt.Println("Null request")
		return
	}

	var requestReply interface{}

	arrivingGuest := decodeRequest(request)
	arrivingGuestName := mux.Vars(request)["name"]

	// Get guest data from guest list
	var guest database.GuestList
	database.Connector.Where("name = ?", arrivingGuestName).Find(&guest)

	// Check if arriving guest is in the checklist
	if guest.Name == "" {
		requestReply = "Guest " + arrivingGuestName + " is not in the guest list"
	} else
	// Check table capacity
	if arrivingGuest.AccompanyingGuests > guest.Table {
		requestReply = "Guest " + arrivingGuestName + " arrived with an entourage bigger than the registered one"
	} else
	// Check if guest already checked in
	if guest.TimeArrived != "" {
		requestReply = "Guest " + arrivingGuestName + " already checked in"
	} else {

		// Update guest data
		guest.AccompanyingGuests = arrivingGuest.AccompanyingGuests
		guest.TimeArrived = utils.GetHoursAndMinutesString()

		// Update guest in the database
		database.Connector.Save(&guest)

		requestReply = CreateCheckInGuestResponse(guest)
	}

	encodeResponse(response, requestReply)
}

// checkOutGuest Processes the request that happens when a guest leaves the party
//
// When a guest leaves, all their accompanying guests leave as well.
func checkOutGuest(response http.ResponseWriter, request *http.Request) {
	if database.Connector == nil {
		fmt.Println("Database unreachable")
		return
	}

	if request == nil {
		fmt.Println("Null request")
		return
	}

	var requestReply interface{}

	// Get guest data from guest list
	var guest database.GuestList
	guestName := mux.Vars(request)["name"]
	database.Connector.Where("name = ?", guestName).Find(&guest)

	// Check if guest is in the checklist
	if guest.Name == "" {
		requestReply = "Guest " + guestName + " is not in the guest list"
	} else
	// Check if guest checked in
	if guest.TimeArrived == "" {
		requestReply = "Guest " + guestName + " has not arrived yet"
	} else {

		// Delete checked in guest from database
		database.Connector.Where("name = ?", guestName).Delete(&guest)

		requestReply = "Guest " + guestName + " left the party"
	}

	encodeResponse(response, requestReply)
}

// getArrivedGuests Processes the request to get the list of guests that have arrived to the party
func getArrivedGuests(response http.ResponseWriter, _ *http.Request) {
	if database.Connector == nil {
		fmt.Println("Database unreachable")
		return
	}

	// Get all guests from database
	var guestList []database.GuestList
	database.Connector.Find(&guestList)

	//Only account for guests that checked in
	guestListIndex := 0
	for _, guest := range guestList {
		if guest.TimeArrived != "" {
			//Copy data and increment index
			guestList[guestListIndex] = guest
			guestListIndex++
		}
	}
	// "Truncate" slice
	guestList = guestList[:guestListIndex]

	encodeResponse(response, CreateGetArrivedGuestsResponse(guestList))
}

// getArrivedGuests Processes the request to get the number of empty seats
func getNumberOfEmptySeats(response http.ResponseWriter, _ *http.Request) {
	if database.Connector == nil {
		fmt.Println("Database unreachable")
		return
	}

	numberOfEmptySeats := 0

	// Get all guests from database
	var guestList []database.GuestList
	database.Connector.Find(&guestList)

	for _, guest := range guestList {

		if guest.TimeArrived == "" {
			// Guest did not check int: empty table
			numberOfEmptySeats += guest.Table
		} else {
			// Guest checked in: check remaining space at table
			numberOfEmptySeats += guest.Table - guest.AccompanyingGuests
		}
	}

	encodeResponse(response, CreateGetNumberOfEmptySeatsResponse(numberOfEmptySeats))
}
