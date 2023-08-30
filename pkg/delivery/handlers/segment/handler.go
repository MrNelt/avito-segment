package segment

import (
	"errors"
	"log"
	"net/http"
	errorType "segment/pkg/errors"
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
	switch {
	case errors.Is(err, errorType.ErrSegmentAlreadyExists):
		log.Println(err.Error())
		errorType.HandleError(ctx, http.StatusBadRequest, err.Error(), err)
		return
	case err != nil:
		log.Println(err.Error())
		errorType.HandleError(ctx, http.StatusInternalServerError, errorType.ErrCreateSegment.Error(), err)
		return
	}
	ctx.Status(http.StatusOK)
}

func (h *Handler) DeleteSegmentByName(ctx *gin.Context) {
	name := ctx.Param("name")
	err := h.repo.DeleteSegmentByName(name)
	if err != nil {
		log.Println(err.Error())
		errorType.HandleError(ctx, http.StatusInternalServerError, errorType.ErrDeleteSegment.Error(), err)
		return
	}
	ctx.Status(http.StatusOK)
}
