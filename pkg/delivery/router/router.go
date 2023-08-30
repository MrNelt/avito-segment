package router

import (
	"net/http"
	"segment/pkg/delivery/handlers/segment"
	"segment/pkg/delivery/handlers/user"
	"segment/pkg/repo"

	"github.com/gin-gonic/gin"
)

func SetupRouter(repo repo.IRepository) *gin.Engine {
	router := gin.Default()
	router.GET("/ping", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "pong")
	})
	var segmentHandler segment.IHandler = segment.NewHandler(repo.GetSegmentRepository())
	var userHandler user.IHandler = user.NewHandler(repo.GetUserRepository())
	segmentGroup := router.Group("/segment")
	{
		segmentGroup.POST("/:name", segmentHandler.CreateSegmentByName)
		segmentGroup.DELETE("/:name", segmentHandler.DeleteSegmentByName)
	}
	userGroup := router.Group("/user")
	{
		userGroup.POST("/:id", userHandler.DeleteAddSegmentsToUser)
		userGroup.GET("/:id", userHandler.GetUserSegmentsByUserID)
	}
	return router
}
