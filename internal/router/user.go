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
	// Protected routes requiring JWT and RBAC

	protectedRouter := r.Group("/").Use(
		middleware.StrictAuth(deps.JWT, deps.Logger),
		middleware.AuthMiddleware(deps.Casbin),
	)
	{
		protectedRouter.GET("/profile", deps.UserHandler.GetProfile)
		protectedRouter.POST("/user/list", deps.UserHandler.GetUserList)
		protectedRouter.PUT("/profile", deps.UserHandler.UpdateProfile)
		protectedRouter.PUT("/user", deps.UserHandler.UpdateUser)
	}
}
