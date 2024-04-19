package common

import (
	"fmt"
	"net/http"
	"out/internal/middleware"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Welcome(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "Welcome to OUT API",
	})
}

func (h *Handler) NotFound(c *gin.Context) {
	e := middleware.NewHttpError("not found", fmt.Sprintf("urlPath:%s", c.Request.URL.Path), http.StatusNotFound)
	c.Error(e)
}
