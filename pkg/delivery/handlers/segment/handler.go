package segment

import (
	"net/http"
	"segment/pkg/errors"
	"segment/pkg/repo/segment"

	"github.com/gin-gonic/gin"
)

type IHandler interface {
	CreateSegmentByName(ctx *gin.Context)
	DeleteSegmentByName(ctx *gin.Context)
}

type Handler struct {
	repo segment.IRepository
}

func NewHandler(repo segment.IRepository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) CreateSegmentByName(ctx *gin.Context) {
	name := ctx.Param("name")
	err := h.repo.CreateSegmentByName(name)
	if err != nil {
		errors.HandleError(ctx, http.StatusInternalServerError, err.Error(), err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"name": name,
	})
}

func (h *Handler) DeleteSegmentByName(ctx *gin.Context) {
	name := ctx.Param("name")
	err := h.repo.DeleteSegmentByName(name)
	if err != nil {
		errors.HandleError(ctx, http.StatusInternalServerError, err.Error(), err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"name": name,
	})
}
