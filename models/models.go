package models

//contact

type Contact struct {
	ID		int32	`json:"id"`
	Title	string	`json:"title"` 
	Completion	bool `json:"isConplete"`
}