package model

import "gorm.io/gorm"

type Permission struct {
	gorm.Model

	// 2. 树形结构支持 (核心)
	// ParentID 用于生成左侧菜单树 (Tree Structure)
	// 0 表示顶级目录
	ParentID int `gorm:"column:parent_id;type:int;default:0" json:"parent_id"`

	// 3. 基础信息
	Name string `gorm:"column:name;type:varchar(50);not null" json:"name"`      // 显示名称: "用户管理"
	Key  string `gorm:"column:key;type:varchar(50);not null;unique" json:"key"` // 唯一标识: "sys:user:list"
	Sort int    `gorm:"column:sort;type:int;default:0" json:"sort"`             // 排序: 数字越小越靠前

	// 4. 类型区分 (非常重要)
	// 1: 目录 (Directory) - 只有名字和图标，不能点击，用于折叠
	// 2: 菜单 (Menu)      - 点击后跳转页面
	// 3: 按钮 (Button/API)- 页面里的按钮，对应后端接口
	Type int `gorm:"column:type;type:tinyint;not null;default:1" json:"type"`

	// 5. 前端路由信息 (Type=2 菜单时必填)
	Path      string `gorm:"column:path;type:varchar(255)" json:"path"`           // 前端路由地址: "/system/user"
	Component string `gorm:"column:component;type:varchar(255)" json:"component"` // 前端组件路径: "views/system/user/index"

	// 6. 后端鉴权信息 (Type=3 按钮时必填)
	// Casbin 或 中间件鉴权时，就匹配这两个字段
	Api    string `gorm:"column:api;type:varchar(255)" json:"api"`      // 接口路径: "/v1/user/:id"
	Method string `gorm:"column:method;type:varchar(10)" json:"method"` // 请求方法: "GET", "POST", "DELETE"
}

func (m *Permission) TableName() string {
	return "permission"
}
