package response

import (
	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func Success(c *gin.Context, code int, status string, message string, data interface{}) {
	c.JSON(code, Response{
		Code:    code,
		Status:  status,
		Message: message,
		Data:    data,
	})
}

func Error(c *gin.Context, code int, status string, message string) {
	c.JSON(code, Response{
		Code:    code,
		Status:  status,
		Message: message,
		Data:    nil,
	})
}

type RegisterMerchantResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Status      string `json:"status"`
	ApiKey      string `json:"api_key"`
	CallbackURL string `json:"callback_url"`
}

type GetMerchantResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Status      string `json:"status"`
	Balance     int64  `json:"balance"`
	CallbackURL string `json:"callback_url"`
}
