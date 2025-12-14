package server

import (
	v1 "go-nunu/api/v1"
	"go-nunu/docs"
	"go-nunu/internal/middleware"
	"go-nunu/internal/router"
	"go-nunu/pkg/server/http"
	"go-nunu/web"
	nethttp "net/http"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewHTTPServer(
	deps router.RouterDeps,
) *http.Server {
	if deps.Config.GetString("env") == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}
	s := http.NewServer(
		gin.Default(),
		deps.Logger,
		http.WithServerHost(deps.Config.GetString("http.host")),
		http.WithServerPort(deps.Config.GetInt("http.port")),
	)

	fsys, err := static.EmbedFolder(web.Assets(), "dist")
	if err != nil {
		panic(err)
	}
	// 设置前端静态资源
	s.Use(static.Serve("/", fsys))
	s.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		if len(path) >= 3 && path[:3] == "/v1" {
			c.JSON(nethttp.StatusNotFound, v1.ErrNotFound)
			return
		}
		indexPageData, err := web.Assets().ReadFile("dist/index.html")
		if err != nil {
			c.String(nethttp.StatusNotFound, "404 page not found")
			return
		}
		c.Data(nethttp.StatusOK, "text/html; charset=utf-8", indexPageData)
	})
	// swagger doc
	docs.SwaggerInfo.BasePath = "/v1"
	s.GET("/swagger/*any", ginSwagger.WrapHandler(
		swaggerfiles.Handler,
		//ginSwagger.URL(fmt.Sprintf("http://localhost:%d/swagger/doc.json", deps.Config.GetInt("app.http.port"))),
		ginSwagger.DefaultModelsExpandDepth(-1),
		ginSwagger.PersistAuthorization(true),
	))

	s.Use(
		middleware.CORSMiddleware(),
		middleware.ResponseLogMiddleware(deps.Logger),
		middleware.RequestLogMiddleware(deps.Logger),
		//middleware.SignMiddleware(log),
		// middleware.AuthMiddleware(deps.Casbin),
	)
	s.GET("/", func(ctx *gin.Context) {
		deps.Logger.WithContext(ctx).Info("hello")
		v1.HandleSuccess(ctx, map[string]interface{}{
			":)": "Thank you for using nunu!",
		})
	})

	v1 := s.Group("/v1")
	router.InitUserRouter(deps, v1)
	router.InitRoleRouter(deps, v1)
	router.InitPermissionRouter(deps, v1)
	router.InitCommonRouter(deps, v1)

	return s
}
