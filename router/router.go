package router

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gcache"
	"github.com/gogf/gf/util/gconv"
	"review-order/app/api/user"
	"review-order/app/pkg/jwt"
	"review-order/app/tools"
	"strconv"
	"time"
)

// 统一路由注册.
func init() {
	s := g.Server()
	s.Group("/api/user", func(g *ghttp.RouterGroup) {
		g.Middleware(MiddlewareAuth)
		g.ALL("/", &user.User{})
	})
	s.Group("/api/menu", func(g *ghttp.RouterGroup) {
		g.Middleware(MiddlewareAuth)
		g.ALL("/", &user.Menu{})
	})
	s.Group("/api/admins", func(g *ghttp.RouterGroup) {
		g.Middleware(MiddlewareAuth)
		g.ALL("/", &user.Admins{})
	})
	s.Group("/api/role", func(g *ghttp.RouterGroup) {
		g.Middleware(MiddlewareAuth)
		g.ALL("/", &user.Role{})
	})
}

func MiddlewareAuth(r *ghttp.Request) {
	var uuid string
	if gconv.String(r.URL) == "/api/user/login" {
		r.Middleware.Next()
	} else {
		token := r.Header.Get(tools.TOKEN_KEY)
		userInfo, ok := jwt.ParseToken(token)
		if !ok {
			tools.ErrorJson(r, "token 无效", 50008)
			return
		}
		exptimestamp, _ := strconv.ParseInt(userInfo["exp"], 10, 64)
		exp := time.Unix(exptimestamp, 0)
		ok = exp.After(time.Now())
		if !ok {
			tools.ErrorJson(r, "token 过期", 50014)
			return
		}
		uuid = userInfo["uuid"]
		if uuid != "" {
			//查询用户ID
			val := gcache.Get(uuid)
			if val == nil {
				tools.ErrorJson(r, "token 无效", 50008)
				return
			}
			userID := gconv.Uint(val)
			r.SetParam(tools.USER_UUID_Key, uuid)
			r.SetParam(tools.USER_ID_Key, userID)
			r.Middleware.Next()
		} else {
			tools.ErrorJson(r, "用户未登录", 50008)
			return
		}

	}

}
