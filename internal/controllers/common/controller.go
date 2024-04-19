package common

import "github.com/gin-gonic/gin"

type Handler struct{}

func RegisterRoutes(r *gin.Engine) {
	h := &Handler{}

	r.GET("/", h.Welcome)
	r.NoRoute(h.NotFound)
}
