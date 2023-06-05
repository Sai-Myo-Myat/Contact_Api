package middlewares

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type response struct {
	ID	int32 `json:"id,omitempty"`
	Message	string	`json:"message,omitempty"`
}

var Schema = `
	CREATE TABLE IF NOT EXISTS contact (
		title text,
		completion BOOLEAN
	)
`

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

func CreateDB() {
	db := createConnection()

	db.MustExec(Schema)
}

func GetAllContacts(w http.ResponseWriter, r *http.Request) {

	res := response{
		ID: 1,
		Message: "Success",
	}
	json.NewEncoder(w).Encode(res)
}