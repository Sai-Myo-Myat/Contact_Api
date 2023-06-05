package middlewares

import (
	"contact-api/models"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type response struct {
	Status int	`json:"status"`
	Message string `json:"message"`
	ID int64 `json:"contact_id"`
	Data []models.Contact  `json:"data"`
}

type contactResponse struct {
	Status int	`json:"status"`
	Message string `json:"message"`
	Data *Contact  `json:"data"`
}

type Contact struct {
	ID int64	`json:"id" db:"id"`
	Name string 	`json:"name" db:"name"`
	PhoneNumber string `json:"phone_number" db:"phone_number"`
	DateOfBirth string `json:"date_of_birth" db:"date_of_birth"`
	Remark string `json:"remark" db:"remark"`
}

func createConnection() *sqlx.DB {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error Loading .env file")
	}

	db, err := sqlx.Connect("postgres", os.Getenv("POSTGRES_URL"))

	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected")

	return db;
}

///MIDDLEWARES

//get all contacts
func GetAllContacts(w http.ResponseWriter, r *http.Request) {

	var contacts []models.Contact

	err := getAllContacts(&contacts)

	var res response

	if err != nil {
		res = response{
			Status: 500,
			Message: err.Error(),
		}
	} else {
		res = response{
			Status: 200,
			Message: "success",
			Data: contacts,
		}
	}

	json.NewEncoder(w).Encode(res)
}

//get one contact
func GetContact(w http.ResponseWriter, r *http.Request){
	var res contactResponse
	var id int64
	vars := mux.Vars(r)
	idString :=vars["id"]
	ctx := r.Context()

	id, err := strconv.ParseInt(idString, 10, 64)

	fmt.Println(id)

	if err != nil {
		res = contactResponse {
			Status: 403,
			Message: err.Error(),
		}
	}

	contact, err := getContact(ctx, id)

	fmt.Println("contact", contact, "error", err, "getContact")

	if err != nil {
		res = contactResponse{
			Status: 404,
			Message: err.Error(),
		}
	}else {
		res = contactResponse{
			Status: 404,
			Message: "Success",
			Data: contact,
		}
	}

	json.NewEncoder(w).Encode(res)
}

//create contact
func CreateContact(w http.ResponseWriter, r *http.Request){
	ctx := r.Context()
	var contact models.Contact

	err := json.NewDecoder(r.Body).Decode(&contact)

	if  err != nil {
		log.Fatalf("Unable to parse request body. %v", err)
	}

	id,  err := insertContact(ctx, contact)

	var res response

	if  err != nil {
		res = response{
			Status: 403,
			Message: err.Error(),
		}
	}else {
		res = response{
			Status: 201,
			Message: "Create contact successfully",
			ID: id,
		}
	}

	json.NewEncoder(w).Encode(res)
}

///DATABASE FUNCTIONS

//get all contacts
func getAllContacts(contacts *[]models.Contact) error {
	db := createConnection()

	queryStatement := `SELECT * FROM contacts`

	err := db.Select(contacts, queryStatement)

	return err
}

//get one contact
func getContact(ctx context.Context, id int64) (*Contact,error){
	db := createConnection()

	queryStatement := `SELECT * FROM contacts WHERE id = $1`

	var rows Contact

	err := db.GetContext(ctx,&rows,queryStatement,id)

	if  err != nil {
		return nil, err
	}

	return &rows, nil
	
}

//insert contact 
func insertContact(ctx context.Context,contact models.Contact) (int64,error) {
	db := createConnection()

	queryStatement := `
	INSERT INTO contacts (name, phone_number, date_of_birth, remark) VALUES (:name, :phone_number, :date_of_birth, :remark)
	RETURNING id`

	stmt, err := db.PrepareNamedContext(ctx, queryStatement)

	var id int64

	err = stmt.GetContext(ctx, &id ,contact)
	if err != nil {
		return 0, err
	}

	return id, nil
}