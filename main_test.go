// main_test.go
package main

import (
	"bytes"

	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
	"github.com/manelmattoussi3/manageCar/connectionDb"
	"github.com/manelmattoussi3/manageCar/structs"
	"github.com/stretchr/testify/assert"

	"github.com/DATA-DOG/go-sqlmock"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var mock sqlmock.Sqlmock
func TestMain(m *testing.M) {
	// Set up a test database
	var err error
	db, mock, err = sqlmock.New()
	if err != nil {
		log.Fatalf("Error creating mock database: %v", err)
	}

	// Initialize the db variable with the mock database
	connectionDb.SetDB(db)

	// Run the tests
	exitCode := m.Run()

	// Close the mock database and exit
	if err := db.Close(); err != nil {
		log.Fatalf("Error closing mock database: %v", err)
	}

	os.Exit(exitCode)
}



func TestGetCars(t *testing.T) {
	// Your test code for GetCars goes here

	// Example test:
	// Create a request and recorder
	req := httptest.NewRequest("GET", "/cars", nil)
	w := httptest.NewRecorder()

	// Create a router and handle the request
	r := setupRouter()
	r.ServeHTTP(w, req)

	// Assert response status code and content
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, w.Code)
	}
}
// Example test function for AddCar
func TestAddCar(t *testing.T) {
	// Create a car to add
	car := structs.Car{
		Model:           "Test Model",
		RegistrationNum: "ABC123",
		Mileage:         100.0,
		Condition:       "available",
	}

	// Serialize the car to JSON
	carJSON, err := json.Marshal(car)
	if err != nil {
		t.Fatal(err)
	}

	// Create a request to test the AddCar function
	req, err := http.NewRequest("POST", "/cars/add", bytes.NewBuffer(carJSON))
	if err != nil {
		t.Fatal(err)
	}

	// Set the request content type to JSON
	req.Header.Set("Content-Type", "application/json")

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Initialize the router and serve the request
	router := setupRouter()
	router.ServeHTTP(rr, req)

	// Check the status code (expecting 200 OK or a suitable status code)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

}
func TestRentCar(t *testing.T) {
	// Initialize your router
	router := mux.NewRouter()
	router.HandleFunc("/cars/{registration}/rentals", RentCar).Methods("POST")

	// Create a request with a registration number
	req, err := http.NewRequest("POST", "/cars/ABC123/rentals", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Serve the request
	router.ServeHTTP(rr, req)

	// Check the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %v, got %v", http.StatusOK, rr.Code)
	}

	// Check the response body if needed
	// Example: Check if the response contains a certain message
	expected := `{"message": "Car rented successfully"}`
	assert.Equal(t, expected, rr.Body.String())
}

func TestReturnCar(t *testing.T) {
	// Initialize your router
	router := mux.NewRouter()
	router.HandleFunc("/cars/{registration}/returns", ReturnCar).Methods("POST")

	// Create a request with a registration number and kilometers driven
	returnData := map[string]float64{"kilometers_driven": 100.0}
	reqBody, _ := json.Marshal(returnData)
	req, err := http.NewRequest("POST", "/cars/ABC123/returns", bytes.NewReader(reqBody))
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Serve the request
	router.ServeHTTP(rr, req)

	// Check the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %v, got %v", http.StatusOK, rr.Code)
	}

	// Check the response body if needed
	// Example: Check if the response contains a certain message
	expected := `{"message": "Car returned successfully"}`
	assert.Equal(t, expected, rr.Body.String())
}
func setupRouter() *mux.Router {
	router := mux.NewRouter()

	// Define your routes and handlers here

	return router
}
