//go:build wireinject
// +build wireinject

package wire

import (
	"go-nunu/internal/repository"
	"go-nunu/internal/server"
	"go-nunu/pkg/app"
	"go-nunu/pkg/log"
	CasbinPkg "go-nunu/pkg/casbin"


	"github.com/google/wire"
	"github.com/spf13/viper"
)

var repositorySet = wire.NewSet(
	repository.NewDB,
	//repository.NewRedis,
	repository.NewRepository,
	repository.NewUserRepository,
	repository.NewRoleRepository,
	repository.NewPermissionRepository,
)
var serverSet = wire.NewSet(
	server.NewMigrateServer,
)

// build App
func newApp(
	migrateServer *server.MigrateServer,
) *app.App {
	return app.NewApp(
		app.WithServer(migrateServer),
		app.WithName("demo-migrate"),
	)
}

// 添加 Casbin 提供者
var casbinSet = wire.NewSet(
    CasbinPkg.NewEnforcer,
)

func NewWire(*viper.Viper, *log.Logger) (*app.App, func(), error) {
	panic(wire.Build(
		repositorySet,
		serverSet,
		casbinSet,
		newApp,
	))
}
