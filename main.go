package main

import (
	"app/appMessaggistica/database"
	"app/appMessaggistica/routers"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Nessun file .env trovato (ok se siamo in produzione)")
	}

	connStr := os.Getenv("DATABASE_URL")
	database.InitDB(connStr)

	r := routers.SetupRouter()
	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	}
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}
	if err := r.Run(port); err != nil {
		log.Fatal(err)
	}
}
