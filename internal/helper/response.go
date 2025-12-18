package helper

import (
	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func SuccessResponse(c *gin.Context, code int, success bool, message string, data any) {
	response := Response{
		Success: success,
		Message: message,
		Data:    data,
	}
	c.JSON(code, response)
}

func ErrorResponse(c *gin.Context, code int, success bool, error string) {
	response := Response{
		Success: success,
		Error:   error,
	}
	c.JSON(code, response)
}
