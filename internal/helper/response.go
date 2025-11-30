package helper

import (
	"github.com/gin-gonic/gin"
)

type Response struct {
	Status  bool   `json:"status"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func SuccessResponse(c *gin.Context, code int, status bool, message string, data any) {
	response := Response{
		Status:  status,
		Message: message,
		Data:    data,
	}
	c.JSON(code, response)
}

func ErrorResponse(c *gin.Context, code int, status bool, error string) {
	response := Response{
		Status: status,
		Error:  error,
	}
	c.JSON(code, response)
}
