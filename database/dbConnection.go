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

	createUsers := `CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY, 
        name VARCHAR(50) UNIQUE NOT NULL, 
        password VARCHAR(255) NOT NULL
    );`

	createMessages := `CREATE TABLE IF NOT EXISTS messages (
        id SERIAL PRIMARY KEY, 
        sent_by VARCHAR(50) NOT NULL, 
        content TEXT NOT NULL, 
        sent_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );`

	if _, err := db.Exec(createUsers); err != nil {
		log.Println("Errore creazione users (o esiste già):", err)
	}
	if _, err := db.Exec(createMessages); err != nil {
		log.Println("Errore creazione messages (o esiste già):", err)
	}
}

func GetDB() *sql.DB {
	return db
}
