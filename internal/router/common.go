package router

import (
	"go-nunu/internal/middleware"

	"github.com/gin-gonic/gin"
)

func InitCommonRouter(
	deps RouterDeps,
	r *gin.RouterGroup,
) {
	// Protected routes requiring JWT and RBAC
	protectedRouter := r.Group("/common").Use(
		middleware.StrictAuth(deps.JWT, deps.Logger),
		middleware.AuthMiddleware(deps.Casbin),
	)
	{
		protectedRouter.POST("/upload", deps.CommonHandler.UploadPresignedUrl)
		protectedRouter.GET("/upload_presigned_url", deps.CommonHandler.UploadPresignedUrl)
	}
}
