package structs

import "github.com/google/uuid"
type Car struct {
	ID              uuid.UUID `json:"id"`
	Model           string    `json:"model"`
	RegistrationNum string    `json:"registration_num"`
	Mileage         float64   `json:"mileage"`
	Condition       string    `json:"condition"` // "available" or "rented"
}
