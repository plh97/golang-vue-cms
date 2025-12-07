package service

import (
	"context"
	"go-nunu/internal/repository"
	"go-nunu/pkg/jwt"
	"go-nunu/pkg/log"
	"go-nunu/pkg/sid"

	"github.com/casbin/casbin/v2"
	"gorm.io/gorm"
)

type Service struct {
	Casbin *casbin.CachedEnforcer
	db     *gorm.DB
	logger *log.Logger
	sid    *sid.Sid
	jwt    *jwt.JWT
	tm     repository.Transaction
}

func NewService(
	db *gorm.DB,
	tm repository.Transaction,
	logger *log.Logger,
	sid *sid.Sid,
	jwt *jwt.JWT,
) *Service {
	return &Service{
		db:     db, // <-- 必须将 DB 实例赋值给 db 字段
		logger: logger,
		sid:    sid,
		jwt:    jwt,
		tm:     tm,
	}
}

// 核心：DB(ctx) 方法的定义，供所有子 Service 继承使用
func (s *Service) DB(ctx context.Context) *gorm.DB {
	// 返回一个绑定了上下文的 DB 实例
	return s.db.WithContext(ctx)
}
