package database

import (
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var db *sql.DB

func InitDB(connStr string) {
	var err error
	db, err = sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal("Errore apertura DB: ", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("DB non raggiungibile: ", err)
	}
	log.Println("Database connesso con successo!")
}

func GetDB() *sql.DB {
	return db
}
