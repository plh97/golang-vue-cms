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
	user, err := h.roleService.GetRoleList(ctx)
	if err != nil {
		fmt.Println("error:", err)
		v1.HandleError(ctx, http.StatusOK, v1.ErrBadRequest, nil)
		return
	}

	v1.HandleSuccess(ctx, user)
}
