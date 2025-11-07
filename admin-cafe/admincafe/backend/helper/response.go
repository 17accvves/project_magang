package helper

import "github.com/gin-gonic/gin"

type Response struct {
    Success bool        `json:"success"`
    Message string      `json:"message"`
    Data    interface{} `json:"data"`
}

func JSON(c *gin.Context, statusCode int, success bool, message string, data interface{}) {
    c.JSON(statusCode, Response{
        Success: success,
        Message: message,
        Data:    data,
    })
}