package main

import (
	"app/appMessaggistica/database"
	"app/appMessaggistica/routers"
	"log"
	"os"
	"strings"
)

func main() {
	connStr := "postgres://gouser:gopass@localhost:5432/godb?sslmode=disable"
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
