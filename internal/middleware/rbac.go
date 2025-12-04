package middleware

import (
	"fmt"

	"github.com/casbin/casbin/v2"
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
	return func(c *gin.Context) {
		// method := c.Request.Method
		e.AddPolicy("admin", "/user", "GET")
		e.AddRoleForUser("zhangsan", "admin")
		Check(e, "zhangsan", "/user", "GET")
		Check(e, "zhangsan", "/user", "POST")
		c.Next()
	}
}
