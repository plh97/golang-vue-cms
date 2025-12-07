package repository

import (
	"context"
	"fmt"
	v1 "go-nunu/api/v1"
	"go-nunu/internal/model"

	"gorm.io/gorm"
)

type RoleRepository interface {
	GetRole(ctx context.Context, id int) (*model.Role, error)
	GetRoleList(ctx context.Context, req v1.GetRoleListRequest) ([]model.Role, error)
	GetRoleCount(ctx context.Context, req v1.GetRoleListRequest) (int, error)
	CreateRole(ctx context.Context, role *model.Role) (*model.Role, error)
	UpdateRole(ctx context.Context, role *model.Role) (*model.Role, error)
	GetPermissionsByIds(ctx context.Context, permissionIDs []uint) ([]model.Permission, error)
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

func (r *roleRepository) GetRole(ctx context.Context, id int) (*model.Role, error) {
	var role model.Role
	err := r.db.WithContext(ctx).First(&role, id).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) Get(ctx context.Context, param v1.GetRoleListRequest) *gorm.DB {
	var roles []model.Role
	db := r.db.WithContext(ctx).Model(&roles)
	if param.PageRequest.CurrentPage > 0 {
		db = db.Scopes(model.Paginate(param.PageRequest))
	}
	if param.Name != "" {
		db = db.Where("name LIKE ?", "%"+param.Name+"%")
	}
	if param.ID != 0 {
		db = db.Where("id = ?", param.ID)
	}
	return db
}

func (r *roleRepository) GetRoleList(ctx context.Context, req v1.GetRoleListRequest) ([]v1.RoleInfo, error) {
	var roles []v1.RoleInfo
	err := r.Get(ctx, req).Find(&roles).Error
	if err != nil {
		return nil, err
	}
	// TODO: use for loop to get all roles' permission from casbin
	for i := 0; i < len(roles); i++ {
		role := &roles[i]
		list, err := r.e.GetAllNamedRoles(role.Name)
		if err != nil {
			return nil, fmt.Errorf("获取角色权限失败: %v", err)
		}
		permissions := make([]v1.RoleInfo, len(list))
		for j, permName := range list {
			permissions[j] = v1.RoleInfo{Role: model.Role{
				Name: permName,
			}}
		}
		role.Permissions = permissions
	}
	return roles, nil
}

func (r *roleRepository) GetRoleCount(ctx context.Context, req v1.GetRoleListRequest) (int, error) {
	var count int64
	db := r.Get(ctx, v1.GetRoleListRequest{Name: req.Name, ID: req.ID})
	err := db.Count(&count).Error
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

func (r *roleRepository) GetPermissionsByIds(ctx context.Context, permissionIDs []uint) ([]model.Permission, error) {
	var permissions []model.Permission
	err := r.db.WithContext(ctx).Where("id IN ?", permissionIDs).Find(&permissions).Error
	if err != nil {
		return nil, err
	}
	return permissions, nil
}
