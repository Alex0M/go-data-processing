package reporting

import (
	"out/internal/db"
	"out/internal/middleware"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	DB         *db.MongoDB
	Collection string
}

func RegisterRoutes(r *gin.Engine, db *db.MongoDB, c string) {
	h := &Handler{
		DB:         db,
		Collection: c,
	}

	routes := r.Group("/stream", middleware.JwtAuthMiddleware())
	routes.GET("/:name", h.GetUniqueClientPerContent)
	routes.GET("/clients", h.GetAllUniqueUsers)
	routes.GET("/clients/:id", h.GetStreamsForUser)
	routes.GET("/geo", h.GetUniqueClientsPerState)
	routes.GET("/devices", h.GetUniqueClientsPerDevice)
}
