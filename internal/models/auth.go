package models

type Login struct {
	Password string `json:"password" binding:"required"`
}
