package service

import (
	"context"
	v1 "go-nunu/api/v1"
	"go-nunu/internal/model"
	"go-nunu/internal/repository"

	"github.com/casbin/casbin/v2"
	"gorm.io/gorm"
)

type RoleService interface {
	GetRole(ctx context.Context, id int64) (*model.Role, error)
	CreateRole(ctx context.Context, req v1.CreateRoleRequest) (*model.Role, error)
	GetRoleList(ctx context.Context, req v1.GetRoleListRequest) ([]model.Role, int, error)
	UpdateRolePermissions(ctx context.Context, roleId int64, permissionIds []uint) error
}

func NewRoleService(
	service *Service,
	roleRepository repository.RoleRepository,
	casbin *casbin.CachedEnforcer,
) RoleService {
	return &roleService{
		Service:        service,
		roleRepository: roleRepository,
		Casbin:         casbin,
	}
}

type roleService struct {
	*Service
	roleRepository repository.RoleRepository
	Casbin         *casbin.CachedEnforcer
}

func (s *roleService) GetRole(ctx context.Context, id int64) (*model.Role, error) {
	return s.roleRepository.GetRole(ctx, id)
}

func (s *roleService) CreateRole(ctx context.Context, req v1.CreateRoleRequest) (*model.Role, error) {
	role := &model.Role{
		Name: req.Name,
		Sid:  req.Key,
		// Status: req.Status,
	}
	role, err := s.roleRepository.CreateRole(ctx, role)
	if err != nil {
		return nil, err
	}
	err = s.UpdateRolePermissions(
		ctx,
		int64(role.ID),
		req.PermissionIds,
	)
	if err != nil {
		return nil, err
	}
	return role, err
}

func (s *roleService) GetRoleList(ctx context.Context, req v1.GetRoleListRequest) ([]model.Role, int, error) {
	roles, err := s.roleRepository.GetRoleList(ctx, req)
	if err != nil {
		return nil, 0, err
	}
	count, err1 := s.roleRepository.GetRoleCount(ctx, req)
	if err1 != nil {
		return roles, 0, err
	}
	return roles, count, nil
}

func (s *roleService) UpdateRolePermissions(ctx context.Context, roleID int64, permissionIDs []uint) error {
	var role model.Role

	// 1. 查找角色
	if err := s.DB(ctx).First(&role, roleID).Error; err != nil {
		return err
	}

	// 2. 查询完整的权限对象（包含 Path 和 Method）
	var permissions []model.Permission
	if err := s.DB(ctx).Find(&permissions, permissionIDs).Error; err != nil {
		return err
	}
	if len(permissionIDs) == 0 {
		permissions = nil
	}

	err := s.DB(ctx).Transaction(func(tx *gorm.DB) error {
		// 3. 更新数据库关联表
		if err := tx.Model(&role).Association("Permissions").Replace(&permissions); err != nil {
			return err
		}

		// 4. 同步 Casbin 策略（关键步骤）
		sub := role.Sid

		// 4.1 清除该角色的旧策略
		if _, err := s.Casbin.RemoveFilteredPolicy(0, sub); err != nil {
			return err
		}

		// 4.2 添加新策略
		var rules [][]string
		for _, perm := range permissions {
			if perm.Path == "" || perm.Method == "" {
				continue
			}
			// 对应 Casbin 规则: p, role_key, path, method
			switch perm.Type {
			case model.PermissionTypeDirectory:
				rules = append(rules, []string{sub, model.DirectoryResourcePrefix + perm.Api, perm.Method})
			case model.PermissionTypeMenu:
				rules = append(rules, []string{sub, model.MenuResourcePrefix + perm.Api, perm.Method})
			case model.PermissionTypeButton:
				rules = append(rules, []string{sub, model.ApiResourcePrefix + perm.Api, perm.Method})
			}
		}

		if len(rules) > 0 {
			if _, err := s.Casbin.AddPolicies(rules); err != nil {
				return err
			}
		}

		return nil
	})

	return err
}
