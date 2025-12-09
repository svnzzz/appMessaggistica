package initializers

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnvVariables() []byte {
	err := godotenv.Load()
	if err != nil {
		fmt.Print("Errore:", err)
	}

	jwtSecret := os.Getenv("JWT_SECRET")

	if jwtSecret == "" {
		log.Fatal("Variabili d'ambiente Mongo mancanti")
	}
	return []byte(jwtSecret)
}
