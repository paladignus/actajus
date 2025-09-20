// Package dto provides the data transfer objects of the application
package dto

type PeopleInputDTO struct {
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	BirthDate     string `json:"birth_date"`
	MotherName    string `json:"mother_name"`
	FatherName    string `json:"father_name"`
	Gender        string `json:"gender"`
	MaritalStatus string `json:"marital_status"`
}
