package response

import "out/internal/models"

func ErrorResponse(e error) models.Response {
	return models.Response{
		Status:  false,
		Message: "Error",
		Data:    e,
	}
}

func SuccessResponse(m string, data any) models.Response {
	return models.Response{
		Status:  true,
		Message: m,
		Data:    data,
	}
}
