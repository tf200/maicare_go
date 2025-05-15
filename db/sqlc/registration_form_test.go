package db

import (
	"testing"

	"github.com/go-faker/faker/v4"
)

func TestCreateRegistrationForm(t *testing.T) {
	arg := CreateRegistrationFormParams{
		ClientFirstName:   faker.FirstName(),
		ClientLastName:    faker.LastName(),
		ClientBsnNumber:   "123456789",
		ClientGender:      "male",
		ClientNationality: "Dutch",
	}

	arg.ClientCity = "Amsterdam"
}
