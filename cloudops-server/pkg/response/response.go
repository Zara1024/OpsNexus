package response

import "github.com/gin-gonic/gin"

// Body 定义统一响应体。
type Body struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// Success 返回成功响应。
func Success(c *gin.Context, data interface{}) {
	c.JSON(200, Body{Code: 200, Message: "success", Data: data})
}

// Fail 返回业务失败响应。
func Fail(c *gin.Context, code int, message string) {
	c.JSON(200, Body{Code: code, Message: message, Data: gin.H{}})
}

// Error 返回HTTP错误响应。
func Error(c *gin.Context, httpCode int, code int, message string) {
	c.JSON(httpCode, Body{Code: code, Message: message, Data: gin.H{}})
}
