package controllers

import (
	"app/appMessaggistica/database"
	"app/appMessaggistica/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func SendMessage(c *gin.Context) {
	var input struct {
		UserID  int    `json:"user_id"`
		Content string `json:"content"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dati non validi"})
		return
	}

	db := database.GetDB()

	var userName string
	err := db.QueryRow("SELECT name FROM users WHERE id = $1", input.UserID).Scan(&userName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Utente non trovato"})
		return
	}

	query := `INSERT INTO messages (sent_by, content, sent_at) VALUES ($1, $2, $3) RETURNING id`

	var newID int
	err = db.QueryRow(query, userName, input.Content, time.Now()).Scan(&newID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Errore salvataggio: " + err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Inviato!", "id": newID})
}

func GetMessages(c *gin.Context) {
	db := database.GetDB()

	rows, err := db.Query("SELECT id, sent_by, content, sent_at FROM messages ORDER BY sent_at ASC LIMIT 50")
	if err != nil {
		c.JSON(500, gin.H{"error": "Errore DB"})
		return
	}
	defer rows.Close()

	var messages []models.Message

	for rows.Next() {
		var m models.Message
		if err := rows.Scan(&m.ID, &m.UserName, &m.Content, &m.CreatedAt); err != nil {
			continue
		}
		messages = append(messages, m)
	}
	c.JSON(200, messages)
}
