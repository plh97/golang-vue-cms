package service

import (
	"context"
	v1 "go-nunu/api/v1"
	"go-nunu/internal/model"
	"go-nunu/internal/repository"
)

type PermissionService interface {
	GetPermissionList(ctx context.Context, req v1.GetPermissionListRequest) ([]model.Permission, int, error)
	CreatePermission(ctx context.Context, req v1.CreatePermissionRequest) (*model.Permission, error)
}

func NewPermissionService(
	service *Service,
	permissionRepository repository.PermissionRepository,
) PermissionService {
	return &permissionService{
		Service:              service,
		permissionRepository: permissionRepository,
	}
}

type permissionService struct {
	*Service
	permissionRepository repository.PermissionRepository
}

func (s *permissionService) GetPermissionList(ctx context.Context, req v1.GetPermissionListRequest) ([]model.Permission, int, error) {

	count, err1 := s.permissionRepository.GetPermissionCount(ctx, req)
	if err1 != nil {
		return nil, 0, err1
	}
	permissions, err2 := s.permissionRepository.GetPermissionList(ctx, req)
	return permissions, count, err2
}

func (s *permissionService) CreatePermission(ctx context.Context, req v1.CreatePermissionRequest) (*model.Permission, error) {
	permission := &model.Permission{
		Name:      req.Name,
		Key:       req.Key,
		Path:      req.Path,
		Method:    req.Method,
		ParentID:  req.ParentID,
		Sort:      req.Sort,
		Type:      req.Type,
		Component: req.Component,
		Api:       req.Api,
	}
	return s.permissionRepository.CreatePermission(ctx, permission)
}
