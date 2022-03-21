package requestRouting

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

// Router Http request router
var Router *mux.Router

// Setup Setups http request Router
//
// Matches incoming requests to their respective handler
func Setup() {

	Router = mux.NewRouter().StrictSlash(true)
	Router.HandleFunc("/guest_list/{name}", addGuest).Methods(http.MethodPost)
	Router.HandleFunc("/guest_list", getGuestList).Methods(http.MethodGet)
	Router.HandleFunc("/guests/{name}", checkInGuest).Methods(http.MethodPut)
	Router.HandleFunc("/guests/{name}", checkOutGuest).Methods(http.MethodDelete)
	Router.HandleFunc("/guests", getArrivedGuests).Methods(http.MethodGet)
	Router.HandleFunc("/seats_empty", getNumberOfEmptySeats).Methods(http.MethodGet)

	fmt.Println("Request Router successfully setup")
}

// ListenForRequests Listens for incoming http requests
func ListenForRequests() {
	routingSetupError := http.ListenAndServe(networkAddress, Router)
	if routingSetupError != nil {
		fmt.Println(routingSetupError.Error())
		panic("Failed to setup request Router")
	}
}
