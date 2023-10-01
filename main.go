package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/google/uuid" // Import the uuid package
	"github.com/manelmattoussi3/manageCar/structs"
	"github.com/manelmattoussi3/manageCar/connectionDb"

)

var cars []structs.Car
var db *sql.DB
func main() {
	// Initialize the database connection
	var err error
	db, err = connectionDb.InitializeDB() // Assign the database connection to the db variable
	if err != nil {
		// Handle the error
		panic(err)
	}
	defer connectionDb.CloseDB(db)
	r := mux.NewRouter()

	// Define API endpoints
	r.HandleFunc("/cars", GetCars).Methods("GET")
	r.HandleFunc("/cars/add", AddCar).Methods("POST")
	// Define API endpoint for renting a car by its registration number.
	r.HandleFunc("/cars/{registration}/rentals", RentCar).Methods("POST")
	// Define API endpoint for returning a car by its registration number.
	r.HandleFunc("/cars/{registration}/returns", ReturnCar).Methods("POST")

	http.Handle("/", r)

	// Start the HTTP server
	http.ListenAndServe(":8087", nil)
}

// GetCars retrieves a list of all cars from the database and returns them as JSON.
func GetCars(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Query all cars from the database
	rows, err := db.Query("SELECT id, model, registration_num, mileage, `condition` FROM car")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Create a slice to store the cars
	var cars []structs.Car

	// Iterate through the rows and scan each car into the slice
	for rows.Next() {
		var car structs.Car
		if err := rows.Scan(&car.ID, &car.Model, &car.RegistrationNum, &car.Mileage, &car.Condition); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		cars = append(cars, car)
	}

	// Check for any errors in row iteration
	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Encode the cars as JSON and send the response
	if err := json.NewEncoder(w).Encode(cars); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}




// AddCar adds a new car to the parking lot and inserts it into the database.
func AddCar(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var car structs.Car
	err := json.NewDecoder(r.Body).Decode(&car)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if a car with the same registration number already exists in the database
	existingCar, err := connectionDb.FindCarByRegistrationNum(db,car.RegistrationNum)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if existingCar != nil {
		http.Error(w, "Car with the same registration number already exists", http.StatusConflict)
		return
	}

	// Generate a new UUID for the car
	car.ID = uuid.New()

	// Insert the new car into the database
	err = connectionDb.InsertCar(db,car)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the inserted ID
	response := map[string]uuid.UUID{"id": car.ID}
	json.NewEncoder(w).Encode(response)
}

func RentCar(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extract the registration number from the URL path parameters.
	vars := mux.Vars(r)
	registrationNum := vars["registration"]

	// Find the car by its registration number in the database.
	rentedCar, err := connectionDb.FindCarByRegistrationNum(db,registrationNum)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if the car exists in the database.
	if rentedCar == nil {
		http.Error(w, "Car not found", http.StatusNotFound)
		return
	}

	// Check if the car is already rented in the database.
	if rentedCar.Condition == "rented" {
		http.Error(w, "Car is already rented", http.StatusConflict)
		return
	}

	// Update the car's status in the database to "rented."
	rentedCar.Condition = "rented"
	err = connectionDb.UpdateCarCondition(db,rentedCar.ID, "rented")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return a success message or the updated car information.
	json.NewEncoder(w).Encode(rentedCar)
}

// ReturnCar returns a car by its registration number.
func ReturnCar(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extract the registration number from the URL path parameters.
	vars := mux.Vars(r)
	registrationNum := vars["registration"]

	// Find the car by its registration number in the database.
	returnedCar, err := connectionDb.FindCarByRegistrationNum(db,registrationNum)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if the car exists.
	if returnedCar == nil {
		http.Error(w, "Car not found", http.StatusNotFound)
		return
	}

	// Check if the car was marked as rented.
	if returnedCar.Condition != "rented" {
		http.Error(w, "Car was not marked as rented", http.StatusConflict)
		return
	}

	// Parse the number of kilometers driven from the request body.
	var returnData struct {
		KilometersDriven float64 `json:"kilometers_driven"`
	}

	err = json.NewDecoder(r.Body).Decode(&returnData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Update the car's mileage and condition in the database.
	returnedCar.Mileage += returnData.KilometersDriven
	returnedCar.Condition = "available"

	// Update the car's information in the database.
	err = connectionDb.UpdateCar(db,returnedCar)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the updated car information.
	json.NewEncoder(w).Encode(returnedCar)
}

