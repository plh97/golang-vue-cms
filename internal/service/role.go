package service

import (
	"context"
	"fmt"
	v1 "go-nunu/api/v1"
	"go-nunu/internal/model"
	"go-nunu/internal/repository"

	"github.com/casbin/casbin/v2"
)

type RoleService interface {
	GetRole(ctx context.Context, id int) (*model.Role, error)
	CreateRole(ctx context.Context, req v1.CreateRoleRequest) (*model.Role, error)
	GetRoleList(ctx context.Context, req v1.GetRoleListRequest) ([]model.Role, int, error)
	UpdateRolePermissions(ctx context.Context, roleId int, permissionIds []uint) error
}

func NewRoleService(
	service *Service,
	roleRepository repository.RoleRepository,
	casbin casbin.CachedEnforcer,
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
	Casbin         casbin.CachedEnforcer
}

func (s *roleService) GetRole(ctx context.Context, id int) (*model.Role, error) {
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
		int(role.ID),
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

func (s *roleService) UpdateRolePermissions(ctx context.Context, roleID int, permissionIDs []uint) error {
	// 1. 必须先获取角色信息，拿到角色的 Key (Sid)
	role, err := s.roleRepository.GetRole(ctx, roleID)
	if err != nil {
		return fmt.Errorf("获取角色失败: %v", err)
	}

	// 2. 获取权限列表
	permissions, err := s.roleRepository.GetPermissionsByIds(ctx, permissionIDs)
	if err != nil {
		return fmt.Errorf("获取权限失败: %v", err)
	}

	// 3. 更新 Casbin 策略
	// Casbin 中的 Subject 应该是角色的 Key (例如 "admin")，而不是 ID
	sub := role.Sid

	// 3.1 先清除该角色在 Casbin 中的旧策略
	// RemoveFilteredPolicy(0, sub) 表示删除第一列(v0/subject)等于 sub 的所有规则
	// if _, err := s.Casbin.Enforcer.RemoveFilteredPolicy(0, sub); err != nil {
	// 	return fmt.Errorf("清除旧策略失败: %v", err)
	// }

	// 3.2 添加新策略
	var rules [][]string
	for _, perm := range permissions {
		if perm.Path == "" || perm.Method == "" {
			continue
		}
		// 对应 Casbin 规则: p, role_key, path, method
		rules = append(rules, []string{sub, perm.Path, perm.Method})
	}

	if len(rules) > 0 {
		if _, err := s.Casbin.Enforcer.AddPolicies(rules); err != nil {
			return fmt.Errorf("添加新策略失败: %v", err)
		}
	}
	return nil
}
