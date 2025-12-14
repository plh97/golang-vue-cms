package server

import (
	"context"
	"go-nunu/internal/model"
	"go-nunu/pkg/log"
	"os"

	"github.com/casbin/casbin/v2"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type MigrateServer struct {
	db     *gorm.DB
	log    *log.Logger
	casbin *casbin.CachedEnforcer
}

func NewMigrateServer(db *gorm.DB, log *log.Logger, casbin *casbin.CachedEnforcer) *MigrateServer {
	return &MigrateServer{
		db:     db,
		log:    log,
		casbin: casbin,
	}
}
func (m *MigrateServer) Start(ctx context.Context) error {
	m.db.Migrator().DropTable(
		&model.User{},
		&model.Role{},
		&model.Permission{},
	)
	if err := m.db.AutoMigrate(
		&model.User{},
		&model.Role{},
		&model.Permission{},
	); err != nil {
		m.log.Error("err: ", zap.Error(err))
		return err
	}
	if err := m.initAdminUser(); err != nil {
		m.log.Error("initAdminUser error", zap.Error(err))
		return err
	}
	if err := m.initRBACAndDemoData(); err != nil {
		m.log.Error("initRBACAndDemoData error", zap.Error(err))
		return err
	}
	m.log.Info("AutoMigrate success")
	os.Exit(0)
	return nil
}
func (m *MigrateServer) Stop(ctx context.Context) error {
	m.log.Info("AutoMigrate stop")
	return nil
}
func (m *MigrateServer) initAdminUser() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	if err != nil {
		m.log.Error("bcrypt.GenerateFromPassword error", zap.Error(err))
		return err
	}
	return m.db.Create(&model.User{
		BaseModel: model.BaseModel{ID: 1},
		UserId:    model.AdminUserID,
		Password:  string(hashedPassword),
		Email:     "admin@gmail.com",
		Name:      "Administrator",
	}).Error
}

// 批量初始化权限、角色、用户、关联关系
func (m *MigrateServer) initRBACAndDemoData() error {
	// 1. 权限定义（菜单和API）
	permissions := []model.Permission{
		// 菜单权限
		{Model: gorm.Model{}, Name: "用户管理", Key: "menu:user", Type: model.PermissionTypeMenu, Path: "/user", Component: "views/user/index"},
		{Model: gorm.Model{}, Name: "角色管理", Key: "menu:role", Type: model.PermissionTypeMenu, Path: "/role", Component: "views/role/index"},
		// API权限
		{Model: gorm.Model{}, Name: "获取用户列表", Key: "api:user:list", Type: model.PermissionTypeButton, Api: "/v1/user/list", Method: "GET"},
		{Model: gorm.Model{}, Name: "创建用户", Key: "api:user:create", Type: model.PermissionTypeButton, Api: "/v1/user", Method: "POST"},
		{Model: gorm.Model{}, Name: "更新用户", Key: "api:user:update", Type: model.PermissionTypeButton, Api: "/v1/user", Method: "PUT"},
		{Model: gorm.Model{}, Name: "获取角色列表", Key: "api:role:list", Type: model.PermissionTypeButton, Api: "/v1/role/list", Method: "GET"},
		{Model: gorm.Model{}, Name: "获取权限列表", Key: "api:permission:list", Type: model.PermissionTypeButton, Api: "/v1/permission/list", Method: "GET"},
		{Model: gorm.Model{}, Name: "上传文件", Key: "api:common:upload", Type: model.PermissionTypeButton, Api: "/v1/common/upload", Method: "POST"},
		{Model: gorm.Model{}, Name: "获取个人信息", Key: "api:profile:get", Type: model.PermissionTypeButton, Api: "/v1/profile", Method: "GET"},
		{Model: gorm.Model{}, Name: "更新个人信息", Key: "api:profile:put", Type: model.PermissionTypeButton, Api: "/v1/profile", Method: "PUT"},
	}
	if err := m.db.Create(&permissions).Error; err != nil {
		return err
	}

	// 2. 角色定义
	adminRole := model.Role{Name: "管理员", Sid: "admin"}
	devRole := model.Role{Name: "开发者", Sid: "dev"}
	if err := m.db.Create(&adminRole).Error; err != nil {
		return err
	}
	if err := m.db.Create(&devRole).Error; err != nil {
		return err
	}

	// 3. 用户定义
	hashPwd := func(p string) string {
		h, _ := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
		return string(h)
	}
	adminUser := model.User{UserId: "admin-uid", Email: "admin@gmail.com", Password: hashPwd("123456"), Name: "Admin"}
	devUser := model.User{UserId: "user-uid", Email: "user@gmail.com", Password: hashPwd("123456"), Name: "DevUser"}
	if err := m.db.Create(&adminUser).Error; err != nil {
		return err
	}
	if err := m.db.Create(&devUser).Error; err != nil {
		return err
	}

	// 4. 角色-权限关联（通过 service 层，自动同步 Casbin）
	var allPerms []model.Permission
	if err := m.db.Find(&allPerms).Error; err != nil {
		return err
	}
	var devPerms []model.Permission
	var uploadPerm *model.Permission
	var updateProfilePerm *model.Permission
	for i, p := range allPerms {
		if p.Key == "api:common:upload" {
			uploadPerm = &allPerms[i]
		}
		if p.Key == "api:profile:put" {
			updateProfilePerm = &allPerms[i]
		}
		if p.Method == "GET" || p.Type == model.PermissionTypeMenu {
			devPerms = append(devPerms, p)
		}
	}
	// dev 也分配上传图片权限
	if uploadPerm != nil {
		devPerms = append(devPerms, *uploadPerm)
	}
	// dev 也分配更新个人信息权限
	if updateProfilePerm != nil {
		devPerms = append(devPerms, *updateProfilePerm)
	}

	// 手动同步 GORM 关联和 Casbin
	if err := m.db.Model(&adminRole).Association("Permissions").Replace(allPerms); err != nil {
		return err
	}
	if err := m.db.Model(&devRole).Association("Permissions").Replace(devPerms); err != nil {
		return err
	}

	// Casbin 同步
	// admin
	m.casbin.RemoveFilteredPolicy(0, adminRole.Sid)
	for _, p := range allPerms {
		if p.Method != "" {
			var obj string
			switch p.Type {
			case model.PermissionTypeMenu:
				obj = "menu:" + p.Path
			case model.PermissionTypeButton:
				obj = "api:" + p.Api
			default:
				continue
			}
			m.casbin.AddPolicy(adminRole.Sid, obj, p.Method)
		}
	}
	// dev
	m.casbin.RemoveFilteredPolicy(0, devRole.Sid)
	for _, p := range devPerms {
		if p.Method != "" {
			var obj string
			switch p.Type {
			case model.PermissionTypeMenu:
				obj = "menu:" + p.Path
			case model.PermissionTypeButton:
				obj = "api:" + p.Api
			default:
				continue
			}
			m.casbin.AddPolicy(devRole.Sid, obj, p.Method)
		}
	}

	// 5. 用户-角色关联
	if err := m.db.Model(&adminUser).Association("Roles").Replace([]model.Role{adminRole}); err != nil {
		return err
	}
	if err := m.db.Model(&devUser).Association("Roles").Replace([]model.Role{devRole}); err != nil {
		return err
	}

	return nil
}
