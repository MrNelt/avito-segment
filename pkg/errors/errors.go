package errors

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

var (
	ErrSegmentAlreadyExists    = fmt.Errorf("ErrSegmentAlreadyExists")
	ErrGetUserSegmentsByUserID = fmt.Errorf("ErrGetUserSegmentsByUserID")
	ErrDeleteAddSegmentsToUser = fmt.Errorf("ErrDeleteAddSegmentsToUser")
	ErrParseDeleteAddSegments  = fmt.Errorf("ErrParseDeleteAddSegments")
	ErrParseIDtoINT            = fmt.Errorf("ErrParseIDtoINT")
	ErrSegmentNotFound         = fmt.Errorf("ErrSegmentNotFound")
	ErrUserNotFound            = fmt.Errorf("ErrUserNotFound")
	ErrCreateSegment           = fmt.Errorf("ErrCreateSegment")
	ErrDeleteSegment           = fmt.Errorf("ErrDeleteSegment")
)

func HandleError(ctx *gin.Context, status int, errMsg string, err error) {
	response := gin.H{
		"error": errMsg,
	}
	ctx.JSON(status, response)
}
