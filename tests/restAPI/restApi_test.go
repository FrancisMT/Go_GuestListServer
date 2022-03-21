package restapitest

import (
	"bytes"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"guestListChallenge/src/database"
	"guestListChallenge/src/requestRouting"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// testCasesRESTAPI Structure holding all the test scenarios
var testCasesRESTAPI = []struct {
	testCaseName     string
	requestPath      string
	requestType      string
	requestContent   map[string]interface{}
	expectedResponse interface{}
}{
	{
		"Adding a valid guest to the guest list",
		"/guest_list/Francisco",
		http.MethodPost,
		map[string]interface{}{
			"table":               5,
			"accompanying_guests": 2,
		},
		requestRouting.CreateAddGuestResponse(
			database.GuestList{
				Name: "Francisco",
			}),
	},
	{
		"Adding a valid guest to the guest list with an entourage bigger than the table capacity",
		"/guest_list/Francisco",
		http.MethodPost,
		map[string]interface{}{
			"table":               1,
			"accompanying_guests": 2},
		"Guest will no be added to the guest list: guest's table cannot hold so many people.",
	},
	{
		"Getting guest list",
		"/guest_list",
		http.MethodGet,
		map[string]interface{}{},
		requestRouting.CreateGetGuestListResponse(
			[]database.GuestList{
				{
					Name:               "Francisco",
					Table:              5,
					AccompanyingGuests: 5,
				},
				{
					Name:               "Martins",
					Table:              4,
					AccompanyingGuests: 2,
				},
			}),
	},
	{
		"Checking in valid guest",
		"/guests/Martins",
		http.MethodPut,
		map[string]interface{}{
			"accompanying_guests": 2,
		},
		requestRouting.CreateCheckInGuestResponse(
			database.GuestList{
				Name: "Martins",
			}),
	},
	{
		"Checking in an invalid guest",
		"/guests/ForeverAlone",
		http.MethodPut,
		map[string]interface{}{
			"accompanying_guests": 0,
		},
		"Guest ForeverAlone is not in the guest list",
	},
	{
		"Checking in a valid guest with an entourage bigger than the table capacity",
		"/guests/Martins",
		http.MethodPut,
		map[string]interface{}{
			"accompanying_guests": 9000,
		},
		"Guest Martins arrived with an entourage bigger than the registered one",
	},
	{
		"Checking in a valid guest that has already checked in",
		"/guests/Francisco",
		http.MethodPut,
		map[string]interface{}{
			"accompanying_guests": 5,
		},
		"Guest Francisco already checked in",
	},
	{
		"Checking out valid guest",
		"/guests/Francisco",
		http.MethodDelete,
		map[string]interface{}{},
		"Guest Francisco left the party",
	},
	{
		"Checking out guest that hasn't checked int yet",
		"/guests/Martins",
		http.MethodDelete,
		map[string]interface{}{},
		"Guest Martins has not arrived yet",
	},
	{
		"Checking out invalid guest",
		"/guests/ForeverAlone",
		http.MethodDelete,
		map[string]interface{}{},
		"Guest ForeverAlone is not in the guest list",
	},
	{
		"Getting list of guests that are already in the party",
		"/guests",
		http.MethodGet,
		map[string]interface{}{},
		requestRouting.CreateGetArrivedGuestsResponse(
			[]database.GuestList{
				{
					Name:               "Francisco",
					AccompanyingGuests: 5,
					TimeArrived:        "13:37",
				},
			}),
	},
	{
		"Getting number of empty seats",
		"/seats_empty",
		http.MethodGet,
		map[string]interface{}{},
		requestRouting.CreateGetNumberOfEmptySeatsResponse(4),
	},
}

// resetDatabase Resets and populates the database with data needed for testing
func resetDatabase() {

	// Delete database contents
	database.Connector.Delete(&database.GuestList{})

	// Populate database
	guests := []database.GuestList{
		{
			Name:               "Francisco",
			Table:              5,
			AccompanyingGuests: 5,
			TimeArrived:        "13:37",
		},
		{
			Name:               "Martins",
			Table:              4,
			AccompanyingGuests: 2,
		},
	}
	for _, guest := range guests {
		database.Connector.Create(&guest)
	}
}

// TestMain Setups all the necessary dependencies for the testing scenarios
func TestMain(m *testing.M) {

	// Create pool for MySQL Docker container
	pool, operationError := dockertest.NewPool("")
	if operationError != nil {
		fmt.Println(operationError.Error())
		panic("Could not connect to Docker")
	}

	// Setup MySQL container options
	opts := dockertest.RunOptions{
		Repository: "mysql",
		Tag:        "5.7",
		Env: []string{
			"MYSQL_ROOT_PASSWORD=password",
			"MYSQL_DATABASE=getground",
			"MYSQL_USER=francisco",
			"MYSQL_PASSWORD=password"},
		ExposedPorts: []string{"3306"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"3306": {
				{HostIP: "127.0.0.1", HostPort: "3306"},
			},
		},
	}

	// Run MySQL Docker container
	resource, operationError := pool.RunWithOptions(&opts)
	if operationError != nil {
		fmt.Println(operationError.Error())
		panic("Could not start MySQL Docker container")
	}

	// Exponential backoff mechanism to wait for MySQL boot
	if operationError := pool.Retry(func() error {

		database.Connector, operationError = gorm.Open("mysql", fmt.Sprintf("francisco:password@(localhost:%s)/getground", resource.GetPort("3306/tcp")))
		if operationError != nil {
			fmt.Println("MySQL database still booting")
			return operationError
		}

		// Check if database is reachable
		return database.Connector.DB().Ping()
	}); operationError != nil {
		fmt.Println(operationError.Error())
		panic("Could not connect to MySQL Docker container")
	}

	// Setup test database
	database.Connector.AutoMigrate(&database.GuestList{})

	// Setup request router
	requestRouting.Setup()

	// Run test scenarios
	code := m.Run()

	// Delete MySQL Docker container
	if operationError := pool.Purge(resource); operationError != nil {
		fmt.Println(operationError.Error())
		panic("Could not purge MySQL Docker container")
	}

	os.Exit(code)
}

// TestRESAPI Runs all the test scenarios
func TestRESAPI(t *testing.T) {
	for _, testCase := range testCasesRESTAPI {

		t.Logf("Running test case: %q\n", testCase.testCaseName)

		resetDatabase()

		// Create request body
		requestBody, err := json.Marshal(testCase.requestContent)
		if err != nil {
			t.Fatalf("Couldn't create request body: %v\n", err)
		}

		// Create request
		request, err := http.NewRequest(testCase.requestType, testCase.requestPath, bytes.NewReader(requestBody))
		if err != nil {
			t.Fatalf("Couldn't create request: %v\n", err)
		}

		// Send request and register response
		responseRecorder := httptest.NewRecorder()
		requestRouting.Router.ServeHTTP(responseRecorder, request)

		// Check response correctness
		if responseRecorder.Code != http.StatusOK {
			t.Errorf("Wrong http status received %d\n", responseRecorder.Code)
		}

		// Prep response for analysis
		receivedResponse := strings.TrimSuffix(responseRecorder.Body.String(), "\n")
		expectedResponse, err := json.Marshal(testCase.expectedResponse)
		if err != nil {
			t.Fatalf("Couldn't encode expected response to Json: %v\n", err)
		}

		// Assert response correctness
		if receivedResponse != string(expectedResponse) {
			t.Errorf("Incorrect response:\nexpected:%q\nreceived:%q \n", string(expectedResponse), receivedResponse)
		} else {
			t.Logf("Expected and received reponses math")
		}
	}
}
