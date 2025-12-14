package router

import (
	"go-nunu/internal/middleware"

	"github.com/gin-gonic/gin"
)

func InitRoleRouter(
	deps RouterDeps,
	r *gin.RouterGroup,
) {
	// Protected routes requiring JWT and RBAC
	protectedRouter := r.Group("/").Use(
		middleware.StrictAuth(deps.JWT, deps.Logger),
		middleware.AuthMiddleware(deps.Casbin),
	)
	{
		protectedRouter.GET("/role/list", deps.RoleHandler.GetRoleList)
		protectedRouter.POST("/role", deps.RoleHandler.CreateRole)
		protectedRouter.PUT("/role", deps.RoleHandler.UpdateRolePermissions)
	}

	// Strict permission routing group
	// strictAuthRouter := r.Group("/").Use(middleware.StrictAuth(deps.JWT, deps.Logger))
	// {
	// }
}
