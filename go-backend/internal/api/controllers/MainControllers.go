package controllers

import "github.com/gin-gonic/gin"

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
	response := CodeResponse{Message: 200}
	context.HTML(200, "index.html", gin.H{"response": response})
}
