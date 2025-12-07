package v1

import (
	"go-nunu/api"
	"go-nunu/internal/model"
)

type CreateRoleRequest struct {
	Name          string `json:"name"`
	Key           string `json:"key"`
	Status        int    `json:"status"`
	PermissionIds []uint `json:"permission_ids"`
}

type GetRoleListRequest struct {
	api.PageRequest
	Name string `json:"name" form:"name"`
	ID   int    `json:"id" form:"id"`
}

type RoleInfo struct {
	model.Role
	Permissions []model.Permission `json:"permissions"`
}

type GetRoleListResponseData struct {
	api.PageResponse
	List []model.Role `json:"list"`
}

type GetRoleListResponse struct {
	Response
	Data GetRoleListResponseData `json:"data"`
}

type UpdateRolePermissionsRequest struct {
	ID            int  `json:"id"`
	PermissionIds []uint `json:"permission_ids"`
}
type CreateRoleResponse struct {
	Response
	Data model.Role `json:"data"`
}
