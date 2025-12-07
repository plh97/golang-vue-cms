package middleware

import (
	"fmt"
	v1 "go-nunu/api/v1"
	"go-nunu/internal/model"
	"go-nunu/pkg/jwt"
	"net/http"

	"github.com/casbin/casbin/v2"
	"github.com/duke-git/lancet/v2/convertor"
	"github.com/gin-gonic/gin"
)

func Check(e *casbin.CachedEnforcer, sub, obj, act string) {
	ok, _ := e.Enforce(sub, obj, act)
	if ok {
		fmt.Printf("%s CAN %s %s\n", sub, act, obj)
	} else {
		fmt.Printf("%s CANNOT %s %s\n", sub, act, obj)
	}
}

func AuthMiddleware(e *casbin.CachedEnforcer) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 从上下文获取用户信息（假设通过 JWT 或其他方式设置）
		v, exists := ctx.Get("claims")
		if !exists {
			v1.HandleError(ctx, http.StatusUnauthorized, v1.ErrUnauthorized, nil)
			ctx.Abort()
			return
		}
		uid := v.(*jwt.MyCustomClaims).UserId
		if convertor.ToString(uid) == model.AdminUserID {
			// 防呆设计，超管跳过API权限检查
			ctx.Next()
			return
		}
		// TODO: 这里可以根据需要添加更多的上下文信息，例如用户角色等
		// 获取请求的资源和操作
		sub := convertor.ToString(uid)
		obj := model.ApiResourcePrefix + ctx.Request.URL.Path
		act := ctx.Request.Method

		// 检查权限
		allowed, err := e.Enforce(sub, obj, act)
		if err != nil {
			v1.HandleError(ctx, http.StatusForbidden, v1.ErrForbidden, nil)
			ctx.Abort()
			return
		}
		if !allowed {
			v1.HandleError(ctx, http.StatusForbidden, v1.ErrForbidden, nil)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
