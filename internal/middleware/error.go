package middleware

import (
	"fmt"
	"net/http"
	"out/internal/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type HttpError struct {
	Message    string `json:"message"`
	Metadata   string `json:"-"`
	StatusCode int    `json:"statusCode"`
}

func (e HttpError) Error() string {
	return fmt.Sprintf("message: %s,  metadata: %s", e.Message, e.Metadata)
}

func NewHttpError(message, metadata string, statusCode int) HttpError {
	return HttpError{
		Message:    message,
		Metadata:   metadata,
		StatusCode: statusCode,
	}
}

func ErrorHandler(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		for _, err := range c.Errors {
			switch e := err.Err.(type) {
			case HttpError:
				logger.Error("API HTTP error", zap.Error(e))
				c.AbortWithStatusJSON(e.StatusCode, response.ErrorResponse(e))
			default:
				logger.Error("internal server error", zap.Error(e))
				c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"message": "Internal Server Error"})
			}
		}
	}
}
