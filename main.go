package main

import (
	"chat.com/route"
	"chat.com/utils"
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"log"
	"os"
)

// 전역변수
var db *sql.DB

func main() {
	utils.HandleErr(godotenv.Load(".env"))

	cfg := mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   os.Getenv("DBHOST") + ":" + os.Getenv("DBPORT"),
		DBName: os.Getenv("DBNAME"),
	}

	// Get a database handle.
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	utils.HandleErr(db.Ping())
	fmt.Println("DB Connected!")

	port := 8080
	route.Start(port)
}
