package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/denisenkom/go-mssqldb"
)

var authenticatedUserID = 2 //todo build a tokenized login/logout mechanism

var db *sql.DB

var server = "localhost"
var port = 1433
var user = "sa"
var password = "sek1gunzC8"
var database = "GocialDB"

func main() {
	connectionString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		server, user, password, port, database)

	var err error

	db, err = sql.Open("sqlserver", connectionString)
	if err != nil {
		log.Fatal("Error creating connection pool: ", err.Error())
	}
	ctx := context.Background()
	err = db.PingContext(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("Connected!\n")

	router := NewRouter()

	log.Fatal(http.ListenAndServe(":8888", router))
}
