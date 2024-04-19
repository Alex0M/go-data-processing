package auth

import (
	"net/http"
	"os"
	"out/internal/middleware"
	"out/internal/models"
	"out/internal/utils"

	"github.com/gin-gonic/gin"
)

var authPass = os.Getenv("AUTH_PASSWORD")

func (h *Handler) Auth(c *gin.Context) {
	var login models.Login

	if err := c.ShouldBindJSON(&login); err != nil {
		e := middleware.NewHttpError("bad request", err.Error(), http.StatusBadRequest)
		c.Error(e)
		return
	}

	if login.Password != authPass {
		e := middleware.NewHttpError("invalid token", "invalid token", http.StatusUnauthorized)
		c.Error(e)
		return
	}

	token, err := utils.GenerateToken(login.Password)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
