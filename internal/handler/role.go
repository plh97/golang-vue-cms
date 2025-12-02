package handler

import (
	"fmt"
	v1 "go-nunu/api/v1"
	"go-nunu/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RoleHandler struct {
	*Handler
	roleService service.RoleService
}

func NewRoleHandler(
	handler *Handler,
	roleService service.RoleService,
) *RoleHandler {
	return &RoleHandler{
		Handler:     handler,
		roleService: roleService,
	}
}

func (h *RoleHandler) GetRoleList(ctx *gin.Context) {
	roleList, err := h.roleService.GetRoleList(ctx)
	if err != nil {
		fmt.Println("error:", err)
		v1.HandleError(ctx, http.StatusOK, v1.ErrBadRequest, nil)
		return
	}

	v1.HandleSuccess(ctx, roleList)
}


func (h *RoleHandler) CreateRole(ctx *gin.Context) {
	var req v1.UpdateProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	role, err := h.roleService.CreateRole(ctx)
	if err != nil {
		fmt.Println("error:", err)
		v1.HandleError(ctx, http.StatusOK, v1.ErrBadRequest, nil)
		return
	}

	v1.HandleSuccess(ctx, role)
}
