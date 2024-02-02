package models

type MakeShorterRequest struct {
	Link string `form:"full_link" json:"full_link" binding:"required"`
}
