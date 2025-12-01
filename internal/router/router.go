package router

import (
	"go-nunu/internal/handler"
	"go-nunu/pkg/jwt"
	"go-nunu/pkg/log"

	"github.com/spf13/viper"
)

type RouterDeps struct {
	Logger            *log.Logger
	Config            *viper.Viper
	JWT               *jwt.JWT
	UserHandler       *handler.UserHandler
	CommonHandler     *handler.CommonHandler
	RoleHandler       *handler.RoleHandler
	PermissionHandler *handler.PermissionHandler
}
