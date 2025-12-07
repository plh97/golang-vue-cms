package handler

import (
	"fmt"
	"go-nunu/api"
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

// GetRoleList godoc
//
//	@Summary	获取用户信息
//	@Schemes
//	@Description
//	@Tags		Role模块
//	@Accept		json
//	@Produce	json
//	@Security	Bearer
//	@Success	200	{object}	v1.GetRoleListRequest
//	@Router		/role/list [get]
func (h *RoleHandler) GetRoleList(ctx *gin.Context) {
	var req v1.GetRoleListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	roleList, count, err := h.roleService.GetRoleList(ctx, req)
	if err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	v1.HandleSuccess(ctx, v1.GetRoleListResponseData{
		List: roleList,
		PageResponse: api.PageResponse{
			Total: count,
		},
	})
}

func (h *RoleHandler) CreateRole(ctx *gin.Context) {
	var req v1.CreateRoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	role, err := h.roleService.CreateRole(ctx, req)
	if err != nil {
		fmt.Println("error:", err)
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	v1.HandleSuccess(ctx, role)
}

// UpdateRolePermissions
func (h *RoleHandler) UpdateRolePermissions(ctx *gin.Context) {
	var req v1.UpdateRolePermissionsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	// Here you would typically call a service method to update the role's permissions.
	// For example:
	err := h.roleService.UpdateRolePermissions(ctx, req.ID, req.PermissionIds)
	if err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	// For now, we'll just return a success response.
	v1.HandleSuccess(ctx, nil)
}
