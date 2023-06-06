package models

//contact

type Contact struct {
	ID int64	`json:"id" db:"id"`
	Name string 	`json:"name" db:"name"`
	PhoneNumber string `json:"phone_number" db:"phone_number"`
	DateOfBirth string `json:"date_of_birth" db:"date_of_birth"`
	Remark string `json:"remark" db:"remark"`
}