package middleware

import (
	"net/http"

	"github.com/RaymondCode/simple-demo/service"
	"github.com/gin-gonic/gin"
)

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

// Auth 鉴权中间件
// 若用户携带的token正确,解析token,将userId放入上下文context中并放行;否则,返回错误信息
func JwtAuth() gin.HandlerFunc {
	return func(context *gin.Context) {
		auth := context.Query("token")
		if len(auth) == 0 {
			context.Abort()
			context.JSON(http.StatusUnauthorized, Response{
				StatusCode: -1,
				StatusMsg:  "Unauthorized",
			})
		}
		claims, flag := service.ParseToken(auth)
		if !flag {
			context.Abort()
			context.JSON(http.StatusUnauthorized, Response{
				StatusCode: -1,
				StatusMsg:  "Token Error",
			})
		} else {
			println("token 正确")
		}
		context.Set("userId", claims.UserId)
		context.Next()
	}
}
