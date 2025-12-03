package repository

import (
	"context"
	v1 "go-nunu/api/v1"
	"go-nunu/internal/model"

	"gorm.io/gorm"
)

type RoleRepository interface {
	GetRole(ctx context.Context, id int64) (*model.Role, error)
	GetRoleList(ctx context.Context, req v1.GetRoleListRequest) ([]model.Role, error)
	GetRoleCount(ctx context.Context, req v1.GetRoleListRequest) (int, error)
	CreateRole(ctx context.Context, role *model.Role) (*model.Role, error)
	UpdateRole(ctx context.Context, role *model.Role) (*model.Role, error)
}

func NewRoleRepository(
	repository *Repository,
) RoleRepository {
	return &roleRepository{
		Repository: repository,
	}
}

type roleRepository struct {
	*Repository
}

func (r *roleRepository) GetRole(ctx context.Context, id int64) (*model.Role, error) {
	var role model.Role

	return &role, nil
}

func (r *roleRepository) Get(ctx context.Context, param v1.GetRoleListRequest) *gorm.DB {
	var roles []model.Role
	db := r.db.WithContext(ctx).Model(&roles)
	if param.PageRequest.CurrentPage > 0 {
		db = db.Scopes(model.SetPage(*param.PageRequest))
	}
	return db
}

func (r *roleRepository) GetRoleList(ctx context.Context, req v1.GetRoleListRequest) ([]model.Role, error) {
	var roles []model.Role
	err := r.Get(ctx, req).Preload("Permissions").Find(&roles).Error
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *roleRepository) GetRoleCount(ctx context.Context, req v1.GetRoleListRequest) (int, error) {
	var count int64
	var roles []model.Role
	err := r.db.WithContext(ctx).Model(&roles).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

// CreateRole 创建角色
func (r *roleRepository) CreateRole(ctx context.Context, role *model.Role) (*model.Role, error) {
	err := r.db.WithContext(ctx).Create(role).Error
	return role, err
}

func (r *roleRepository) UpdateRole(ctx context.Context, role *model.Role) (*model.Role, error) {
	err := r.db.WithContext(ctx).Save(role).Error
	return role, err
}
