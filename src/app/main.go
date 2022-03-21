package main

import (
	"guestListChallenge/src/database"
	"guestListChallenge/src/requestRouting"
)

// main App entrypoint
func main() {
	database.Connect()
	requestRouting.Setup()
	requestRouting.ListenForRequests()
}
