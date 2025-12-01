package model

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	Name        string       `gorm:"column:name;type:varchar(50);not null;unique" json:"name"`
	Key         string       `gorm:"column:key;type:varchar(50);not null;unique" json:"key"`
	Status      int          `gorm:"column:status;type:tinyint;default:1" json:"status"`
	Permissions []Permission `gorm:"many2many:sys_role_permissions;" json:"permissions"`
}

func (m *Role) TableName() string {
	return "role"
}
