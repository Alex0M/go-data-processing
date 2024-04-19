package tracking

import (
	"encoding/json"
	"fmt"
	"net/http"
	"out/internal/middleware"
	"out/internal/models"
	"out/internal/response"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
)

func (h *Handler) ReportStream(c *gin.Context) {
	var stream models.Stream
	if err := c.ShouldBindJSON(&stream); err != nil {
		e := middleware.NewHttpError("bad request", err.Error(), http.StatusBadRequest)
		c.Error(e)
		return
	}

	s, err := json.Marshal(stream)
	if err != nil {
		e := middleware.NewHttpError("bad request", err.Error(), http.StatusBadRequest)
		c.Error(e)
		return
	}

	partition, offset, err := h.KP.SendMessage(&sarama.ProducerMessage{
		Topic: h.Topic,
		//Key:   sarama.StringEncoder(stream.ClientID),
		Value: sarama.StringEncoder(s),
	})

	if err != nil {
		e := middleware.NewHttpError("failed to store your data", err.Error(), http.StatusInternalServerError)
		c.Error(e)
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Current Unique Users streaming content", fmt.Sprintf("Your data is stored with unique identifier important/%d/%d", partition, offset)))
}
