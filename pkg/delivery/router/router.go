package router

import (
	"net/http"
	"segment/pkg/delivery/handlers/segment"
	"segment/pkg/repo"

	"github.com/gin-gonic/gin"
)

func SetupRouter(repo repo.IRepository) *gin.Engine {
	router := gin.Default()
	router.GET("/ping", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "pong")
	})
	var segmentHandler segment.IHandler = segment.NewHandler(repo.GetSegmentRepository())
	router.POST("/segment/:name", segmentHandler.CreateSegmentByName)
	router.DELETE("/segment/:name", segmentHandler.DeleteSegmentByName)
	return router
}
