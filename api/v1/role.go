package v1

import "go-nunu/internal/model"

type CreateRoleRequest struct {
	Name          string `json:"name"`
	Key           string `json:"key"`
	Status        int    `json:"status"`
	PermissionIds []uint `json:"permission_ids"`
}

type GetRoleListRequest struct {
	PageRequest
}

type GetRoleListResponseData struct {
	PageResponse
	List []model.Role `json:"list"`
}

type GetRoleListResponse struct {
	Response
	Data GetRoleListResponseData `json:"data"`
}

type UpdateRolePermissionsRequest struct {
	ID            int64  `json:"id"`
	PermissionIds []uint `json:"permission_ids"`
}
type CreateRoleResponse struct {
	Response
	Data model.Role `json:"data"`
}
