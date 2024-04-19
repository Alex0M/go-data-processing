package reporting

import (
	"fmt"
	"net/http"
	"out/internal/response"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetAllUniqueUsers(c *gin.Context) {
	cc, err := h.DB.GetAllUniqueUsers(c.Request.Context(), h.Collection)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Current Unique users across all streams", cc))
}

func (h *Handler) GetStreamsForUser(c *gin.Context) {
	userID := c.Param("id")
	cc, err := h.DB.GetStreamsForUser(c.Request.Context(), h.Collection, userID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse(fmt.Sprintf("Current stream counts for %s", userID), cc))
}

func (h *Handler) GetUniqueClientPerContent(c *gin.Context) {
	streamName := c.Param("name")
	cc, err := h.DB.GetUniqueClientPerContent(c.Request.Context(), h.Collection, streamName)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse(fmt.Sprintf("Current Unique Users streaming content on %s", streamName), cc))
}

func (h *Handler) GetUniqueClientsPerState(c *gin.Context) {
	cc, err := h.DB.GetUniqueClientsPerState(c.Request.Context(), h.Collection)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Current Unique users across all streams per state", cc))
}

func (h *Handler) GetUniqueClientsPerDevice(c *gin.Context) {
	cc, err := h.DB.GetUniqueClientsPerDevice(c.Request.Context(), h.Collection)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Current Unique users across all streams per device", cc))
}
