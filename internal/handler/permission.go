package handler

import (
	"fmt"
	"go-nunu/api"
	v1 "go-nunu/api/v1"
	"go-nunu/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PermissionHandler struct {
	*Handler
	permissionService service.PermissionService
}

func NewPermissionHandler(
	handler *Handler,
	permissionService service.PermissionService,
) *PermissionHandler {
	return &PermissionHandler{
		Handler:           handler,
		permissionService: permissionService,
	}
}

func (h *PermissionHandler) GetPermissionList(ctx *gin.Context) {
	// Implementation for getting a permission
	var req v1.GetPermissionListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	permissionList, count, err := h.permissionService.GetPermissionList(ctx, req)
	if err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	v1.HandleSuccess(ctx, v1.GetPermissionListResponseData{
		List: permissionList,
		PageResponse: api.PageResponse{
			Total: count,
		},
	})
}

func (h *PermissionHandler) CreatePermission(ctx *gin.Context) {
	var req v1.CreatePermissionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	permission, err := h.permissionService.CreatePermission(ctx, req)
	if err != nil {
		fmt.Println("error:", err)
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	v1.HandleSuccess(ctx, permission)
}
