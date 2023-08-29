package errors

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

var (
	ErrSegmentNotFound = fmt.Errorf("segment doesn't exist")
	ErrUserNotFound    = fmt.Errorf("user doesn't exist")
)

func HandleError(ctx *gin.Context, status int, errMsg string, err error) {
	response := gin.H{
		"error": errMsg,
	}
	if err != nil {
		response["message"] = err.Error()
	}
	ctx.JSON(status, response)
}
