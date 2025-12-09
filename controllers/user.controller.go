package controllers

import (
	"app/appMessaggistica/database"
	"app/appMessaggistica/initializers"
	"app/appMessaggistica/models"
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = initializers.LoadEnvVariables()

func Register(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := database.GetDB()

	var count int
	checkQuery := `SELECT COUNT(*) FROM users WHERE name = $1`

	err := db.QueryRowContext(ctx, checkQuery, user.Name).Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Errore Database check: " + err.Error()})
		return
	}

	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nome utente già esistente"})
		return
	}

	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Errore hashing"})
		return
	}

	insertQuery := `INSERT INTO users (name, password) VALUES ($1, $2) RETURNING id`

	var newID int

	err = db.QueryRowContext(ctx, insertQuery, user.Name, hashedPassword).Scan(&newID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Errore nella registrazione: " + err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"user_id": newID})
}

func Login(c *gin.Context) {
	var input struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var user models.User

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dati non validi"})
		return
	}

	log.Printf("INPUT NAME RICEVUTO: '%s'", input.Name)

	db := database.GetDB()

	checkQuery := `SELECT id, name, password FROM users WHERE name = $1`

	err := db.QueryRowContext(ctx, checkQuery, input.Name).Scan(&user.ID, &user.Name, &user.Password)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Nome utente o password errati"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Errore Database check: " + err.Error()})
		return
	}

	if err = CheckPassword(user.Password, input.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nome o password errati"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		c.JSON(500, gin.H{"error": "Errore nella generazione del token"})
		return
	}

	c.JSON(200, gin.H{
		"token": tokenString,
		"user":  gin.H{"id": user.ID, "nome": user.Name},
	})

}
