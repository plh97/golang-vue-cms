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
	noStrictAuthRouter := r.Group("/").Use(middleware.NoStrictAuth(deps.JWT, deps.Logger))
	{
		noStrictAuthRouter.POST("/role/list", deps.RoleHandler.GetRoleList)
		noStrictAuthRouter.PUT("/role", deps.UserHandler.UpdateProfile)
	}

	// Strict permission routing group
	// strictAuthRouter := r.Group("/").Use(middleware.StrictAuth(deps.JWT, deps.Logger))
	// {
	// }
}
