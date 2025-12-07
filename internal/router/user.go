package router

import (
	"go-nunu/internal/middleware"

	"github.com/gin-gonic/gin"
)

func InitUserRouter(
	deps RouterDeps,
	r *gin.RouterGroup,
) {
	// No route group has permission
	noAuthRouter := r.Group("/")
	{
		noAuthRouter.POST("/register", deps.UserHandler.Register)
		noAuthRouter.POST("/login", deps.UserHandler.Login)
	}
	// Non-strict permission routing group

	noStrictAuthRouter := r.Group("/").Use(
		middleware.NoStrictAuth(deps.JWT, deps.Logger),
		middleware.AuthMiddleware(deps.Casbin),
	)
	{
		noStrictAuthRouter.GET("/profile", deps.UserHandler.GetProfile)
		noStrictAuthRouter.POST("/user/list", deps.UserHandler.GetUserList)
		noStrictAuthRouter.PUT("/profile", deps.UserHandler.UpdateProfile)
		noStrictAuthRouter.PUT("/user", deps.UserHandler.UpdateUser)
	}

	// Strict permission routing group
	// strictAuthRouter := r.Group("/").Use(middleware.StrictAuth(deps.JWT, deps.Logger))
	// {
	// }
}
