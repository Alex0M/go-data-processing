package auth

import "github.com/gin-gonic/gin"

type Handler struct{}

func RegisterRoutes(r *gin.Engine) {
	h := &Handler{}

	routes := r.Group("/auth")
	routes.POST("/login", h.Auth)
}
