package router

import (
	"go-nunu/internal/middleware"

	"github.com/gin-gonic/gin"
)

func InitCommonRouter(
	deps RouterDeps,
	r *gin.RouterGroup,
) {
	// Non-strict permission routing group
	noStrictAuthRouter := r.Group("/common").Use(
		middleware.NoStrictAuth(deps.JWT, deps.Logger),
		middleware.AuthMiddleware(deps.Casbin),
	)
	{
		noStrictAuthRouter.POST("/upload", deps.CommonHandler.UploadPresignedUrl)
		noStrictAuthRouter.GET("/upload_presigned_url", deps.CommonHandler.UploadPresignedUrl)
	}
}
