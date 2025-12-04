package router

import (
	"go-nunu/internal/handler"
	"go-nunu/pkg/jwt"
	"go-nunu/pkg/log"

	"github.com/casbin/casbin/v2"
	"github.com/spf13/viper"
)

type RouterDeps struct {
	Logger            *log.Logger
	Config            *viper.Viper
	JWT               *jwt.JWT
	Casbin            *casbin.CachedEnforcer
	UserHandler       *handler.UserHandler
	CommonHandler     *handler.CommonHandler
	RoleHandler       *handler.RoleHandler
	PermissionHandler *handler.PermissionHandler
}
