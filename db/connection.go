package db

import (
	"chat.com/utils"
	"context"
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"log"
	"os"
)

// 전역변수
var (
	ctx context.Context
	db  *sql.DB
)

func Start() {
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
}

func Create(query string, args ...any) {
	fmt.Println(1)
	tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	utils.HandleErr(err)

	fmt.Println(2)
	_, execErr := tx.ExecContext(ctx, query, args...)
	if execErr != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Fatalf("insert failed: %v, unable to rollback: %v\n", execErr, rollbackErr)
		}
		log.Fatalf("insert failed: %v", execErr)
	}
	utils.HandleErr(tx.Commit())
}
