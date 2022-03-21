package database

// dbConnectionConfig Database connection configuration
var dbConnectionConfig = struct {
	User           string
	Password       string
	ServerProtocol string
	ServerName     string
	ServerPort     string
	DBName         string
}{
	User:           "francisco",
	Password:       "password",
	ServerProtocol: "tcp",
	ServerName:     "mysql",
	ServerPort:     "3306",
	DBName:         "getground",
}

// getConnectionString Returns the connection string for the current database setup
func getConnectionString() string {
	const leftParenthesis = "("
	const rightParenthesis = ")"
	const atSymbol = "@"
	const colon = ":"
	const slash = "/"

	return dbConnectionConfig.User + colon +
		dbConnectionConfig.Password + atSymbol +
		dbConnectionConfig.ServerProtocol + leftParenthesis +
		dbConnectionConfig.ServerName + colon +
		dbConnectionConfig.ServerPort + rightParenthesis + slash +
		dbConnectionConfig.DBName

}
