package models

//contact

type Contact struct {
	ID int64 `json:"id"`
	Name 	string 	`json:"name"`
	Phone_Number string `json:"phone_number"`
	Date_Of_Birth string	`json:"date_of_birth"`
	Remark 	string	`json:"remark"`
}