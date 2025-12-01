package service

import (
    "context"
	"go-nunu/internal/model"
	"go-nunu/internal/repository"
)

type RoleService interface {
	GetRole(ctx context.Context, id int64) (*model.Role, error)
	GetRoleList(ctx context.Context) ([]model.Role, error)
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


func (s *roleService) GetRoleList(ctx context.Context) ([]model.Role, error) {
	return s.roleRepository.GetRoleList(ctx)
}
