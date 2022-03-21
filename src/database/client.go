package database

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

//Connector Database connection  for CRUD operation's
var Connector *gorm.DB

//Connect Creates MySQL connection and performs database migration for GuestList
func Connect() {

	var connectionError error

	// Establish connection to the database
	Connector, connectionError = gorm.Open(dbConnectionConfig.ServerName, getConnectionString())
	if connectionError != nil {
		fmt.Println(connectionError.Error())
		panic("Failed to connect to Database")
	}

	fmt.Println("Connection to database was successful")

	// Auto migrate to keep tables reflecting structs
	Connector.AutoMigrate(&GuestList{})
}
