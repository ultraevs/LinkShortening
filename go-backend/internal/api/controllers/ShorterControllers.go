package controllers

import (
	"app/internal/database"
	"app/internal/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"os"
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
	shorted := MakeShort()
	err := SaveLinkToDatebase(request.Link, shorted)
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

func SaveLinkToDatebase(fullLink string, shortLink string) error {
	insertSQL := "INSERT INTO short_link (link, short) VALUES ($1, $2);"
	_, err := database.Db.Exec(insertSQL, fullLink, shortLink)
	if err != nil {
		fmt.Println("Error inserting into short_link:", err)
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

func MakeShort() string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

	b := make([]byte, 5)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
