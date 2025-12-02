package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	UserId   string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	Email    string `gorm:"not null"`
	// CreatedAt time.Time
	// UpdatedAt time.Time
	// DeletedAt gorm.DeletedAt `gorm:"index"`

	// ID              int                   `gorm:"column:id;type:int;primaryKey;autoIncrement:true" json:"id"`                 // 主键ID
	Name            string `gorm:"column:name;type:varchar(255);not null" json:"name"`                      // 用户姓名
	Image           string `gorm:"column:image;type:varchar(255);not null" json:"image"`                    // image
	Gender          int    `gorm:"column:gender;type:tinyint;not null;default 0" json:"gender"`             // 性别 1.男 2.女
	Udid            string `gorm:"column:udid;type:varchar(255);not null" json:"udid"`                      // 设备唯一标识
	UserType        int    `gorm:"column:user_type;type:tinyint;not null" json:"user_type"`                 // 用户类型 1访客用户 2注册用户 3假用户
	IsSubmitProfile int    `gorm:"column:is_submit_profile;type:tinyint;not null" json:"is_submit_profile"` // 是否已经上传过个人资料 1.是 2.否
	Status          int    `gorm:"column:status;type:tinyint;not null" json:"status"`                       // 状态 0未知状态 1正常 2禁用 3注销
	CountryID       int    `gorm:"column:country_id;type:int;not null" json:"country_id"`                   // 国家ID
	RegisterIP      string `gorm:"column:register_ip;type:varchar(45);not null" json:"register_ip"`         // 注册IP
	RegisterTime    int    `gorm:"column:register_time;type:int;not null" json:"register_time"`             // 注册时间
	RegisterType    int    `gorm:"column:register_type;type:tinyint;not null" json:"register_type"`         // 注册类型 0游客 1微信服务号 2微信小程序
	LastLoginIP     string `gorm:"column:last_login_ip;type:varchar(45);not null" json:"last_login_ip"`     // 最后登录IP
	LastLoginTime   int    `gorm:"column:last_login_time;type:int;not null" json:"last_login_time"`         // 最后登录时间
	LastLoginType   int    `gorm:"column:last_login_type;type:tinyint;not null" json:"last_login_type"`     // 最后登录类型 0游客 1微信服务号 2微信小程序
	DeactivateTime  int    `gorm:"column:deactivate_time;type:int;not null" json:"deactivate_time"`         // 注销时间
	Roles           []Role `gorm:"many2many:sys_user_roles;"`                                               // role
}

func (u *User) TableName() string {
	return "users"
}
