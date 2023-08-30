package user

import (
	"errors"
	"log"
	"net/http"
	errorType "segment/pkg/errors"
	"segment/pkg/repo/user"
	"strconv"

	"github.com/gin-gonic/gin"
)

type IHandler interface {
	GetUserSegmentsByUserID(ctx *gin.Context)
	DeleteAddSegmentsToUser(ctx *gin.Context)
}

type Handler struct {
	repo user.IRepository
}

func NewHandler(repo user.IRepository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) DeleteAddSegmentsToUser(ctx *gin.Context) {
	type Request struct {
		Add    []string `json:"add"`
		Delete []string `json:"delete"`
	}
	var req Request
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Println(err.Error())
		errorType.HandleError(ctx, http.StatusBadRequest, errorType.ErrParseDeleteAddSegments.Error(), err)
		return
	}
	IDstr := ctx.Param("id")
	ID, err := strconv.Atoi(IDstr)
	if err != nil {
		log.Println(err.Error())
		errorType.HandleError(ctx, http.StatusBadRequest, errorType.ErrParseIDtoINT.Error(), err)
		return
	}
	err = h.repo.DeleteAddSegmentsToUser(req.Delete, req.Add, ID)
	switch {
	case errors.Is(err, errorType.ErrSegmentNotFound):
		log.Println(err.Error())
		errorType.HandleError(ctx, http.StatusBadRequest, err.Error(), err)
		return
	case err != nil:
		log.Println(err.Error())
		errorType.HandleError(ctx, http.StatusInternalServerError, errorType.ErrDeleteAddSegmentsToUser.Error(), err)
		return
	}
	ctx.Status(http.StatusOK)
}

func (h *Handler) GetUserSegmentsByUserID(ctx *gin.Context) {
	IDstr := ctx.Param("id")
	ID, err := strconv.Atoi(IDstr)
	if err != nil {
		log.Println(err.Error())
		errorType.HandleError(ctx, http.StatusBadRequest, errorType.ErrParseIDtoINT.Error(), err)
		return
	}
	userWithSegments, err := h.repo.GetUserSegmentsByUserID(ID)
	if err != nil {
		log.Println(err.Error())
		errorType.HandleError(ctx, http.StatusInternalServerError, errorType.ErrGetUserSegmentsByUserID.Error(), err)
		return
	}
	ctx.JSON(http.StatusOK, userWithSegments)
}
