package controllers

import (
	"app/internal/database"
	"app/internal/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"os"
	"strings"
)

type Response struct {
	Message string `json:"message"`
}

// ShortRedirect godoc
// @Summary Short Link Redirect
// @Description Redirecting from short link to full
// @Param key path string true "short link key"
// @Produce json
// @Tags Shorter
// @Success 200 {object} Response
// @Router /short/{key} [get]
func ShortRedirect(context *gin.Context) {
	shortLink := context.Param("key")
	fullLink, err := CheckShortLink(shortLink)
	if err != nil {
		response := Response{Message: err.Error()}
		context.JSON(http.StatusInternalServerError, gin.H{"response": response})
		return
	}
	context.Redirect(http.StatusFound, "https://"+fullLink)
	context.Abort()
}

// DoShorter godoc
// @Summary Do Link Short
// @Description Getting full link and make it short
// @Param full_link formData string true "full link"
// @Produce json
// @Tags Shorter
// @Success 200 {object} Response
// @Router /shorter [post]
func DoShorter(context *gin.Context) {
	var request models.MakeShorterRequest
	if err := context.ShouldBind(&request); err != nil {
		response := Response{Message: "error 4xx"}
		context.JSON(http.StatusBadRequest, gin.H{"response": response})
		return
	}
	// Извлекаем значение куки
	cookie, err := context.Request.Cookie("session")
	if err != nil {
		response := Response{Message: "error getting user cookie"}
		context.JSON(http.StatusBadRequest, gin.H{"response": response})
		return
	}
	shorted := MakeShort()
	err = SaveLinkToDatabase(strings.TrimPrefix(request.Link, "https://"), shorted, cookie)
	if err != nil {
		response := Response{Message: err.Error()}
		context.JSON(http.StatusInternalServerError, gin.H{"response": response})
		return
	}
	var resp string
	if os.Getenv("SETUP_TYPE") == "local" {
		resp = "http://localhost:8082/short/" + string(shorted)
	} else {
		resp = "https://shorter.ultraevs.ru/short/" + string(shorted)
	}
	response := Response{Message: resp}
	context.JSON(http.StatusOK, gin.H{"response": response})
}

func SaveLinkToDatabase(fullLink, shortLink string, cookie *http.Cookie) error {
	insertSQL := "INSERT INTO short_link (link, short) VALUES ($1, $2);"
	_, err := database.Db.Exec(insertSQL, fullLink, shortLink)
	if err != nil {
		fmt.Println("Error inserting into short_link:", err)
		return err
	}

	historyUpdateSQL := `
        INSERT INTO shorting_history (cookie, history)
        VALUES ($1, $2::jsonb)
        ON CONFLICT (cookie) DO UPDATE
        SET history = shorting_history.history || $2::jsonb;
    `
	var preShort string
	if os.Getenv("SETUP_TYPE") == "local" {
		preShort = "http://localhost:8082/short/"
	} else {
		preShort = "https://shorter.ultraevs.ru/short/"
	}
	jsonData, err := json.Marshal(map[string]string{"https://" + fullLink: preShort + shortLink})
	if err != nil {
		fmt.Println("Error marshaling history data:", err)
		return err
	}

	_, err = database.Db.Exec(historyUpdateSQL, cookie.Value, jsonData)
	if err != nil {
		fmt.Println("Error updating shorting_history:", err)
		return err
	}

	return err
}

func CheckShortLink(shortLink string) (string, error) {
	var fullLink string
	row := database.Db.QueryRow("SELECT link FROM short_link WHERE short = $1;", shortLink)
	err := row.Scan(&fullLink)
	if err != nil {
		return "error", err
	}
	return fullLink, nil

}

func GetHistoryFromDB(cookie *http.Cookie) (map[string]string, error) {
	historyJSON := ""
	row := database.Db.QueryRow("SELECT history FROM shorting_history WHERE cookie = $1;", cookie.Value)
	err := row.Scan(&historyJSON)
	if err != nil {
		if err == sql.ErrNoRows {
			return make(map[string]string), nil
		}
		return nil, err
	}

	var historyMap map[string]string
	if err := json.Unmarshal([]byte(historyJSON), &historyMap); err != nil {
		return nil, err
	}

	return historyMap, nil
}

func MakeShort() string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

	b := make([]byte, 5)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
