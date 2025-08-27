package web

import "github.com/gin-gonic/gin"

type WebResposne struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func ResponseJSON(ctx *gin.Context, code int, status string, message string, data any) {
	ctx.JSON(code, WebResposne{
		Code:    code,
		Status:  status,
		Message: message,
		Data:    data,
	})
}
