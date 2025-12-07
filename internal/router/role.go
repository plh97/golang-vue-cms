package router

import (
	"go-nunu/internal/middleware"

	"github.com/gin-gonic/gin"
)

func InitRoleRouter(
	deps RouterDeps,
	r *gin.RouterGroup,
) {
	// Non-strict permission routing group
	noStrictAuthRouter := r.Group("/").Use(
		middleware.NoStrictAuth(deps.JWT, deps.Logger),
		middleware.AuthMiddleware(deps.Casbin),
	)
	{
		noStrictAuthRouter.GET("/role/list", deps.RoleHandler.GetRoleList)
		noStrictAuthRouter.POST("/role", deps.RoleHandler.CreateRole)
		noStrictAuthRouter.PUT("/role", deps.RoleHandler.UpdateRolePermissions)
	}

	// Strict permission routing group
	// strictAuthRouter := r.Group("/").Use(middleware.StrictAuth(deps.JWT, deps.Logger))
	// {
	// }
}
