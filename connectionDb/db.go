// db.go
package connectionDb

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/manelmattoussi3/manageCar/structs"
)
var db *sql.DB
// InitializeDB initializes the database connection and returns a *sql.DB instance.
func InitializeDB() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/car_rental")
	if err != nil {
		return nil, err
	}
	return db, nil
}

// CloseDB closes the database connection.
func CloseDB(db *sql.DB) {
	db.Close()
}
// SetDB sets the global database connection.
func SetDB(database *sql.DB) {
	db = database
}
// FindCarByRegistrationNum retrieves a car by its registration number from the database.
func FindCarByRegistrationNum(db *sql.DB,registrationNum string) (*structs.Car, error) {
	var car structs.Car
	query := "SELECT * FROM car WHERE registration_num = ?"

	err := db.QueryRow(query, registrationNum).Scan(&car.ID, &car.Model, &car.RegistrationNum, &car.Mileage, &car.Condition)
	if err == sql.ErrNoRows {
		return nil, nil // Car not found
	} else if err != nil {
		return nil, err // Database error
	}

	return &car, nil
}

// InsertCar inserts a car record into the database.
func InsertCar(db *sql.DB,car structs.Car) error {
	insertSQL := "INSERT INTO car (id, model, registration_num, mileage, `condition`) VALUES (?, ?, ?, ?, ?)"

	_, err := db.Exec(insertSQL, car.ID, car.Model, car.RegistrationNum, car.Mileage, car.Condition)
	if err != nil {
		return err
	}

	return nil
}
// UpdateCarCondition updates the condition of a car in the database.
func UpdateCarCondition(db *sql.DB,id uuid.UUID, condition string) error {
	updateSQL := "UPDATE car SET `condition` = ? WHERE id = ?"

	_, err := db.Exec(updateSQL, condition, id)
	if err != nil {
		return err
	}

	return nil
}
// UpdateCar updates a car's information in the database.
func UpdateCar(db *sql.DB,car *structs.Car) error {
	updateSQL := "UPDATE car SET mileage = ?,`condition` = ? WHERE id = ?"
	_, err := db.Exec(updateSQL, car.Mileage, car.Condition, car.ID)
	if err != nil {
		return err
	}

	return nil
}

