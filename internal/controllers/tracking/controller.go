package tracking

import (
	"out/internal/middleware"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	KP    sarama.SyncProducer
	Topic string
}

func RegisterRoutes(r *gin.Engine, p sarama.SyncProducer, t string) {
	h := &Handler{
		KP:    p,
		Topic: t,
	}

	routes := r.Group("/stream", middleware.JwtAuthMiddleware())
	routes.POST("", h.ReportStream)
}
