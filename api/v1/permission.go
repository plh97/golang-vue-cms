package v1

import (
	"go-nunu/api"
	"go-nunu/internal/model"
)

type GetPermissionListRequest struct {
	api.PageRequest
	ID   int    `json:"id" form:"id"`
	Name string `json:"name" form:"name"`
}

type GetPermissionListResponseData struct {
	api.PageResponse
	List []model.Permission `json:"list"`
}
type GetPermissionListResponse struct {
	Response
	Data GetPermissionListResponseData `json:"data"`
}

type CreatePermissionRequest struct {
	Name      string `json:"name"`
	Key       string `json:"key"`
	// Status    int    `json:"status"`
	Path      string `json:"path"`   // 前端路由地址: "/system/user"
	Method    string `json:"method"` // 请求方法: "GET", "POST", "DELETE"
	ParentID  int    `json:"parent_id"`
	Sort      int    `json:"sort"` // 排序: 数字越小越靠前
	Type      int    `json:"type"`
	Component string `json:"component"` // 前端组件路径: "views/system/user/index"
	Api       string `json:"api"`       // 接口路径: "/v1/user/:id"
}
type CreatePermissionResponse struct {
	Response
	Data model.Permission `json:"data"`
}
