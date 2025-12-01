//go:build wireinject
// +build wireinject

package wire

import (
	"go-nunu/internal/handler"
	"go-nunu/internal/job"
	"go-nunu/internal/repository"
	"go-nunu/internal/router"
	"go-nunu/internal/server"
	"go-nunu/internal/service"
	"go-nunu/pkg/app"
	"go-nunu/pkg/jwt"
	"go-nunu/pkg/log"
	"go-nunu/pkg/server/http"
	"go-nunu/pkg/sid"
	"go-nunu/pkg/aws"

	"github.com/google/wire"
	"github.com/spf13/viper"
)

var repositorySet = wire.NewSet(
	repository.NewDB,
	//repository.NewRedis,
	//repository.NewMongo,
	repository.NewRepository,
	repository.NewTransaction,
	repository.NewUserRepository,
	repository.NewRoleRepository,
	repository.NewPermissionRepository,
)

var serviceSet = wire.NewSet(
	service.NewService,
	service.NewUserService,
	service.NewRoleService,
	service.NewPermissionService,
	service.NewCommonService,
)

var handlerSet = wire.NewSet(
	handler.NewHandler,
	handler.NewUserHandler,
	handler.NewRoleHandler,
	handler.NewPermissionHandler,
	handler.NewCommonHandler,
)

var jobSet = wire.NewSet(
	job.NewJob,
	job.NewUserJob,
)
var serverSet = wire.NewSet(
	server.NewHTTPServer,
	server.NewJobServer,
)

// build App
func newApp(
	httpServer *http.Server,
	jobServer *server.JobServer,
	// task *server.Task,
) *app.App {
	return app.NewApp(
		app.WithServer(httpServer, jobServer),
		app.WithName("demo-server"),
	)
}

// 声明 R2 构造函数
var awsSet = wire.NewSet(
	aws.NewR2Client, // Wire 会自动处理 *viper.Viper 和 *log.Logger 的注入
)

func NewWire(cfg *viper.Viper, logger *log.Logger) (*app.App, func(), error) {
	panic(wire.Build(
		repositorySet,
		serviceSet,
		handlerSet,
		jobSet,
		serverSet,
		wire.Struct(new(router.RouterDeps), "*"),
		sid.NewSid,
		jwt.NewJwt,
		awsSet,
		newApp,
	))
}
