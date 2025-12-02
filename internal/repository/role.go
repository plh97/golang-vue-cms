package repository

import (
	"context"
	"go-nunu/internal/model"
)

type RoleRepository interface {
	GetRole(ctx context.Context, id int64) (*model.Role, error)
	CreateRole(ctx context.Context) (*model.Role, error)
	GetRoleList(ctx context.Context) ([]model.Role, error)
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

func (r *roleRepository) GetRoleList(ctx context.Context) ([]model.Role, error) {
	// 1. 变量名用复数 roles (好习惯)
	var roles []model.Role

	// 2. 链式调用
	// WithContext(ctx): 传递上下文
	// Order("sort desc"): 加上排序（通常角色列表需要按顺序展示）
	// Find(&roles): 查询
	// .Error: 获取错误信息
	err := r.db.WithContext(ctx).Find(&roles).Error

	// 3. 严谨的错误处理
	if err != nil {
		return nil, err
	}

	// 4. 直接返回 slice，不需要取地址 &
	return roles, nil
}


func (r *roleRepository) CreateRole(ctx context.Context) (*model.Role, error) {
	// 1. 变量名用复数 roles (好习惯)
	var role *model.Role

	// 2. 链式调用
	// WithContext(ctx): 传递上下文
	// Order("sort desc"): 加上排序（通常角色列表需要按顺序展示）
	// Find(&roles): 查询
	// .Error: 获取错误信息
	err := r.db.WithContext(ctx).Find(&role).Error

	// 3. 严谨的错误处理
	if err != nil {
		return nil, err
	}

	// 4. 直接返回 slice，不需要取地址 &
	return role, nil
}
