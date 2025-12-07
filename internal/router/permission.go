package router

import (
	"go-nunu/internal/middleware"

	"github.com/gin-gonic/gin"
)

func InitPermissionRouter(
	deps RouterDeps,
	r *gin.RouterGroup,
) {
	// Non-strict permission routing group
	noStrictAuthRouter := r.Group("/").Use(
		middleware.NoStrictAuth(deps.JWT, deps.Logger),
		middleware.AuthMiddleware(deps.Casbin),
	)
	{
		noStrictAuthRouter.GET("/permission/list", deps.PermissionHandler.GetPermissionList)
		noStrictAuthRouter.POST("/permission", deps.PermissionHandler.CreatePermission)
	}

	// Strict permission routing group
	// strictAuthRouter := r.Group("/").Use(middleware.StrictAuth(deps.JWT, deps.Logger))
	// {
	// }
}
