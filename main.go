package main

import (
	"chat.com/route"
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"log"
	"os"
)

var db *sql.DB

// https://go.dev/doc/tutorial/database-access#multiple_rows
func main() {

	fmt.Println("USER : ", os.Getenv("DBUSER"))
	fmt.Println("PASSWORD : ", os.Getenv("DBPASS"))

	cfg := mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "recordings",
	}
	// Get a database handle.
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")

	port := 8080
	route.Start(port)
}
