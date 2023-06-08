package middlewares

import (
	"contact-api/models"
	"context"
	"database/sql"
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
	Data []*models.Contact  `json:"data"`
	Total int64 `json:"total"`
}

type contactResponse struct {
	Status int	`json:"status"`
	Message string `json:"message"`
	Data *models.Contact  `json:"data"`
}

type contactList struct {
	Status int `json:"status"`
	Message string `json:"message"`
	Contacts []*models.Contact `json:"contacts"`
	Total int64 `json:"total"`
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

	return db;
}

///MIDDLEWARES

//get all contacts
func GetAllContacts(w http.ResponseWriter, r *http.Request) {

	var res contactList
	ctx := r.Context();
	query := r.URL.Query()
	limit := int64(6)
	limitString := query.Get("limit")
	if limitString != "" {
		l, err := strconv.ParseInt(limitString, 10, 64)

		if err != nil {
			w.WriteHeader(400)
			res = contactList {
				Status: 400,
				Message: err.Error(),
			}
		}
		limit = l
	}
	

	offset := int64(0);
	offsetString := query.Get("offset");
	if offsetString != "" {
		off, err := strconv.ParseInt(offsetString, 10, 64)

		if err != nil {
			w.WriteHeader(400)
			res = contactList {
				Status: 400,
				Message: err.Error(),
			}
		}
	 	offset= off
	}

	search := query.Get("search");

	contacts,total,err := getAllContacts(ctx, limit, offset, search)

	if err != nil {
		w.WriteHeader(500)
		res = contactList{
			Status: 500,
			Message: err.Error(),
		}
	} else {
		res = contactList{
			Status: 200,
			Message: "success",
			Contacts: contacts,
			Total: total,
		}
	}

	fmt.Println("Get all data", limit, offset)

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

	if err != nil {
		w.WriteHeader(400)
		res = contactResponse{
			Status: 400,
			Message: err.Error(),
		}
	}

	contact, err := getContact(ctx, id)

	// fmt.Println("contact", contact, "error", err, "getContact")
	if err != nil {
		w.WriteHeader(404)
		res = contactResponse{
			Status: 404,
			Message: err.Error(),
		}
	}else {
		res = contactResponse{
			Status: 200,
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

	fmt.Println("body", r.Body)

	err := json.NewDecoder(r.Body).Decode(&contact)

	if  err != nil {
		log.Fatalf("Unable to parse request body. %v", err)
	}

	id,  err := insertContact(ctx, contact)

	var res response

	if  err != nil {
		w.WriteHeader(400)
		res = response{
			Status: 400,
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

//update contact
func UpdateContacts(w http.ResponseWriter, r *http.Request){
	var res response
	idString := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idString, 10, 64);

	if err != nil {
		w.WriteHeader(400)
		res = response{
			Status: 400,
			Message: err.Error(),
		}
	}

	ctx := r.Context()
	var contact models.Contact

	json.NewDecoder(r.Body).Decode(&contact)

	contact.ID = id;

	if err != nil {
		w.WriteHeader(400)
		res = response{
			Status: 400,
			Message: err.Error(),
		}
	}

	err = updateContact(ctx, contact)

	if err != nil {
		w.WriteHeader(404)
		res = response{
			Status: 404,
			Message: err.Error(),
		}
	}else {
		res = response{
			Status: 200,
			Message: "Success",
			ID: id,
		}
	}

	json.NewEncoder(w).Encode(res)
}

func DeleteContact(w http.ResponseWriter, r * http.Request){
	ctx := r.Context();
	idString := mux.Vars(r)["id"];

	id, err := strconv.ParseInt(idString, 10, 64);
	var res response;
	if err != nil {
		w.WriteHeader(400)
		res = response{
			Status: 400,
			Message: err.Error(),
		}
	}

	result, err := deleteContact(ctx, id);


	if  err != nil {
		w.WriteHeader(404)
		res = response{
			Status: 404,
			Message: err.Error(),
		}
	}

	rowsAffected, err	:= result.RowsAffected()

	if  err != nil {
		w.WriteHeader(404)
		res = response{
			Status: 404,
			Message: err.Error(),
		}
	}else if rowsAffected == 0{
		w.WriteHeader(404)
		res = response{
			Status: 404,
			Message: "There is no contact with this id",
		}
	}else {
		res = response{
			Status: 200,
			Message: "Deleted contact successfully",
			ID: id,
		}
	}

	json.NewEncoder(w).Encode(res);
}

///DATABASE FUNCTIONS

//get all contacts
func getAllContacts(ctx context.Context, limit, offset int64, search string) ([]*models.Contact,int64,error) {
	db := createConnection()

	defer db.Close()
	args := map[string]any{
		"limit": limit,
		"offset": offset,
		"name": search,
	}

	// query := ``

	// if search != "" {
	// 	query = `WHERE name :name`
	// }
	queryStatement := `SELECT * FROM contacts ORDER BY id DESC LIMIT :limit OFFSET :offset`

	stmt, err := db.PrepareNamedContext(ctx,queryStatement)

	if err != nil {
		return nil,0, err;
	}

	var contacts []*models.Contact;


	err = stmt.SelectContext(ctx, &contacts, args)
	if err != nil {
		return nil,0,err
	}

	getToTalContactQuery := `SELECT COUNT(*) FROM contacts`

	totalContactStmt, err := db.PrepareNamedContext(ctx, getToTalContactQuery);
	if err != nil {
		return nil, 0, err
	}

	var total int64;

	err = totalContactStmt.GetContext(ctx, &total, args)

	return contacts, total,  nil;
}

//get one contact
func getContact(ctx context.Context, id int64) (*models.Contact,error){
	db := createConnection()

	queryStatement := `SELECT * FROM contacts WHERE id = $1`

	var contact models.Contact

	err := db.GetContext(ctx,&contact,queryStatement,id)

	if  err != nil {
		return nil, err
	}
 
	return &contact, nil
	
}

//insert contact 
func insertContact(ctx context.Context, contact models.Contact) (int64, error) {
	db := createConnection()

	defer db.Close()

	queryStatement := `
	INSERT INTO contacts (name, phone_number, date_of_birth, remark) VALUES (:name, :phone_number, :date_of_birth, :remark)
	RETURNING id`

	stmt, err := db.PrepareNamedContext(ctx, queryStatement)

	var id int64

	err = stmt.GetContext(ctx, &id ,contact)
	fmt.Println(contact.Name, "date")
	if err != nil {
		return 0, err
	}

	return id, nil
}

//update contact 
func updateContact(ctx context.Context, contact models.Contact)  error{
	db := createConnection()

	defer db.Close()
	
	queryStatement := `
		UPDATE contacts SET name = :name, phone_number = :phone_number,
		date_of_birth = :date_of_birth, remark = :remark WHERE 	id = :id;
	`
	_, err := db.NamedExecContext(ctx, queryStatement, contact)


	return err;
}

//delete contact 
func deleteContact(ctx context.Context, id int64) (sql.Result, error) {
	db := createConnection();

	defer db.Close()

	queryStatement := `DELETE FROM contacts WHERE id = $1`;

	result, err := db.ExecContext(ctx, queryStatement, id); 

	if err != nil {
		return nil, err;
	}

	return result, nil;
}