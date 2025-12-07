package model

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	Name string `json:"name" gorm:"column:name;type:varchar(100);uniqueIndex;comment:角色名"`
	Sid  string `json:"sid" gorm:"column:sid;type:varchar(100);uniqueIndex;comment:角色标识"`
	// Key    string `gorm:"column:key;type:varchar(50);not null;unique" json:"key"`
	// Status int    `gorm:"column:status;type:tinyint;default:1" json:"status"`
	// Permissions []Permission `json:"permissions"`
}

func (m *Role) TableName() string {
	return "roles"
}
