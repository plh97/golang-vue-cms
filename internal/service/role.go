package service

import (
	"context"
	v1 "go-nunu/api/v1"
	"go-nunu/internal/model"
	"go-nunu/internal/repository"

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
) RoleService {
	return &roleService{
		Service:        service,
		roleRepository: roleRepository,
	}
}

type roleService struct {
	*Service
	roleRepository repository.RoleRepository
}

func (s *roleService) GetRole(ctx context.Context, id int64) (*model.Role, error) {
	return s.roleRepository.GetRole(ctx, id)
}

func (s *roleService) CreateRole(ctx context.Context, req v1.CreateRoleRequest) (*model.Role, error) {
	role := &model.Role{
		Name:   req.Name,
		Key:    req.Key,
		Status: req.Status,
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

	// 1. 查找角色 (必须先查出来才能操作关联)
	// 注意：ID 统一使用 uint 类型
	if err := s.DB(ctx).First(&role, roleID).Error; err != nil {
		return err // 角色不存在，返回 gorm.ErrRecordNotFound
	}

	var permissions []model.Permission
	for _, id := range permissionIDs {
		permissions = append(permissions, model.Permission{ID: uint(id)})
	}
	err := s.DB(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&role).Association("Permissions").Replace(&permissions); err != nil {
			return err // 关联失败，回滚
		}

		// 4. 返回 nil 提交事务
		return nil
	})

	return err
}
