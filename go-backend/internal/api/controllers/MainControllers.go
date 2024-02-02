package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type CodeResponse struct {
	Message int `json:"code"`
}

// GetMain godoc
// @Summary Get Main
// @Description Get Main
// @Produce json
// @Tags main
// @Success 200 {object} CodeResponse
// @Router / [get]
func GetMain(context *gin.Context) {
	response := Response{Message: "2xx"}
	context.HTML(http.StatusOK, "index.html", gin.H{"response": response})
}

// GetHistory godoc
// @Summary Get History
// @Description Get history of shorting operation
// @Produce json
// @Tags main
// @Success 200 {object} CodeResponse
// @Router /get [get]
func GetHistory(context *gin.Context) {
	// Извлекаем значение куки
	cookie, err := context.Request.Cookie("session")
	if err != nil {
		response := Response{Message: "error getting user cookie"}
		context.JSON(http.StatusBadRequest, gin.H{"response": response})
		return
	}
	dict, err := GetHistoryFromDB(cookie)
	if err != nil {
		response := Response{Message: err.Error()}
		context.JSON(http.StatusInternalServerError, gin.H{"response": response})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": dict})

}
