package service

import (
	"context"
	v1 "go-nunu/api/v1"
	"go-nunu/internal/model"
	"go-nunu/internal/repository"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService interface {
	Register(ctx context.Context, req *v1.RegisterRequest) error
	Login(ctx context.Context, req *v1.LoginRequest) (string, error)
	GetProfile(ctx context.Context, userId string) (*v1.GetProfileResponseData, error)
	GetUserList(ctx context.Context) (*v1.GetUserListResponseData, error)
	UpdateProfile(ctx context.Context, userId string, req *v1.UpdateProfileRequest) error
	UpdateUser(ctx context.Context, req *v1.UpdateUserRequest) error
}

func NewUserService(
	service *Service,
	userRepo repository.UserRepository,
) UserService {
	return &userService{
		userRepo: userRepo,
		Service:  service,
	}
}

type userService struct {
	userRepo repository.UserRepository
	*Service
}

func (s *userService) Register(ctx context.Context, req *v1.RegisterRequest) error {
	// check username
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return v1.ErrInternalServerError
	}
	if err == nil && user != nil {
		return v1.ErrEmailAlreadyUse
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	// Generate user ID
	userId, err := s.sid.GenString()
	if err != nil {
		return err
	}
	user = &model.User{
		UserId:   userId,
		Email:    req.Email,
		Password: string(hashedPassword),
	}
	// Transaction demo
	err = s.tm.Transaction(ctx, func(ctx context.Context) error {
		// Create a user
		if err = s.userRepo.Create(ctx, user); err != nil {
			return err
		}
		// TODO: other repo
		return nil
	})
	return err
}

func (s *userService) Login(ctx context.Context, req *v1.LoginRequest) (string, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil || user == nil {
		return "", v1.ErrUnauthorized
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return "", err
	}
	token, err := s.jwt.GenToken(user.UserId, time.Now().Add(time.Hour*24*90))
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *userService) GetProfile(ctx context.Context, userId string) (*v1.GetProfileResponseData, error) {
	user, err := s.userRepo.GetByID(ctx, userId)
	if err != nil {
		return nil, err
	}

	return &v1.GetProfileResponseData{
		UserId:   user.UserId,
		Nickname: user.Name,
		Email:    user.Email,
		Image:    user.Image,
	}, nil
}

func (s *userService) GetUserList(ctx context.Context) (*v1.GetUserListResponseData, error) {
	user, err := s.userRepo.GetUserList(ctx)
	if err != nil {
		return nil, err
	}
	// count, err1 := s.userRepo.GetUserCount(ctx)
	// if err1 != nil {
	// 	return nil, err1
	// }

	return &v1.GetUserListResponseData{
		List: *user,
		// Total: count,
	}, nil
}

func (s *userService) UpdateProfile(ctx context.Context, userId string, req *v1.UpdateProfileRequest) error {
	user, err := s.userRepo.GetByID(ctx, userId)
	if err != nil {
		return err
	}

	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Nickname != "" {
		user.Name = req.Nickname
	}
	user.Image = req.Image

	if err = s.userRepo.Update(ctx, user); err != nil {
		return err
	}

	return nil
}

func (s *userService) UpdateUser(ctx context.Context, req *v1.UpdateUserRequest) error {
	user, err := s.userRepo.GetByID(ctx, req.UserId)
	if err != nil {
		return err
	}

	var roles []model.Role
	for _, id := range req.RoleIds {
		roles = append(roles, model.Role{ID: uint(id)})
	}

	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Nickname != "" {
		user.Name = req.Nickname
	}
	if req.Image != "" {
		user.Image = req.Image
	}
	err = s.DB(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&user).Association("Roles").Replace(&roles); err != nil {
			return err
		}
		if err = s.userRepo.Update(ctx, user); err != nil {
			return err
		}
		// 4. 返回 nil 提交事务
		return nil
	})
	return err
}
